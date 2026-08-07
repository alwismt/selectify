package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"alwis.dev/selectify/internal/httpx"

	"alwis.dev/selectify/internal/httpx/request"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
	"alwis.dev/selectify/internal/types"
)

const passwordResetTTL = 5 * time.Minute

var (
	ErrUserDeactivated   = errors.New("user is deactivated")
	ErrInvalidResetToken = errors.New("This password reset link is invalid or has expired. Please request a new one.")
)

type AuthService interface {
	RegisterUser(ctx context.Context, req *request.UserRegisterRequest, userAgent, ip string) (*model.UserSession, error)
	LoginUser(ctx context.Context, req *request.LoginRequest, userAgent, ip string) (*model.UserSession, error)
	UserLogout(ctx context.Context, sesId uuid.UUID) (err error)
	ForgetPassword(ctx context.Context, email, ip, userAgent string) error
	ValidateResetToken(ctx context.Context, token string) error
	ResetPassword(ctx context.Context, token, newPassword, ip, userAgent string) error
}

type authService struct {
	jwtService     JWTService
	eventPublisher EventPublisher

	txRepo            repo.TransactionRepo
	userRepo          repo.UserRepo
	userRoleRepo      repo.UserRoleRepo
	sessionRepo       repo.UserSessionRepo
	passwordResetRepo repo.PasswordResetRepo
}

func NewAuthService(jwtService JWTService, eventPublisher EventPublisher, userRepo repo.UserRepo, txRepo repo.TransactionRepo,
	userRole repo.UserRoleRepo, sessionRepo repo.UserSessionRepo, passwordResetRepo repo.PasswordResetRepo) AuthService {
	return &authService{
		jwtService:     jwtService,
		eventPublisher: eventPublisher,

		txRepo:            txRepo,
		userRepo:          userRepo,
		userRoleRepo:      userRole,
		sessionRepo:       sessionRepo,
		passwordResetRepo: passwordResetRepo,
	}
}

func (svc *authService) RegisterUser(ctx context.Context, req *request.UserRegisterRequest, userAgent, ip string) (*model.UserSession, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, logger.Error(ctx, err, "Error while hashing password")
	}

	user := &model.User{
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Status:       types.UserStatusActive,
		PasswordHash: string(hash),
	}

	tx, err := svc.txRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, logger.Error(ctx, err, "Error while begin transaction")
	}

	defer func() {
		if !tx.CanCommit {
			tx.End()
		}
	}()

	if err = svc.userRepo.AddUserWithTx(ctx, tx.Transaction, user); err != nil {
		return nil, logger.Error(ctx, err, "Error while adding user")
	}

	role := &model.UserRole{
		UserID: user.ID,
		Role:   types.RoleCustomer,
	}

	if err = svc.userRoleRepo.InsertUserRoleForCustomerWithTx(ctx, tx.Transaction, role); err != nil {
		return nil, logger.Error(ctx, err, "Error while inserting user role")
	}

	tx.CanCommit = true
	tx.End()

	return svc.newUserSession(ctx, user.ID, userAgent, ip)
}

func (svc *authService) LoginUser(ctx context.Context, req *request.LoginRequest, userAgent, ip string) (*model.UserSession, error) {
	user, err := svc.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, logger.Errorf(ctx, err, "Error while fetching user by email %s", req.Email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		_ = logger.Errorf(ctx, err, "Error while comparing password")
		return nil, httpx.ErrUserNotFound
	}

	if user.Status != types.UserStatusActive {
		return nil, ErrUserDeactivated
	}

	session, err := svc.newUserSession(ctx, user.ID, userAgent, ip)
	if err != nil {
		return nil, err
	}

	if pubErr := svc.eventPublisher.Publish(ctx, types.EventTypeUserLogin, map[string]any{
		"user_id":    user.ID,
		"session_id": session.SessionId.String(),
		"ip":         ip,
		"user_agent": userAgent,
	}); pubErr != nil {
		_ = logger.Error(ctx, pubErr, "failed to publish user.logged_in event")
	}

	return session, nil
}

func (svc *authService) newUserSession(ctx context.Context, userId uint, userAgent, ip string) (*model.UserSession, error) {
	now := time.Now().UTC()
	session := &model.UserSession{
		UserId:    userId,
		UserAgent: userAgent,
		IpAddress: ip,
		SessionId: uuid.New(),
		IssuedAt:  now,
		ExpiresAt: now.Add(20 * time.Minute),
	}

	if err := svc.sessionRepo.InsertUserSession(ctx, session); err != nil {
		return nil, logger.Error(ctx, err, "Error while inserting user session")
	}

	return session, nil
}

func (svc *authService) UserLogout(ctx context.Context, sesId uuid.UUID) (err error) {
	err = svc.sessionRepo.RevokeSession(ctx, sesId)
	if err != nil {
		err = logger.Error(ctx, err, "Error while revoking session")
	}
	return
}

func (svc *authService) ForgetPassword(ctx context.Context, email, ip, userAgent string) error {
	user, err := svc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, httpx.ErrUserNotFound) {
			_ = logger.Errorf(ctx, err, "User not found for email %s", email)
			return nil
		}
		return logger.Errorf(ctx, err, "Error while fetching user by email %s", email)
	}

	if user.Status != types.UserStatusActive {
		_ = logger.Errorf(ctx, ErrUserDeactivated, "Skipping password reset for inactive user %d", user.ID)
		return nil
	}

	rawToken, tokenHash, err := generatePasswordResetToken()
	if err != nil {
		return logger.Error(ctx, err, "Error while generating password reset token")
	}

	if err = svc.passwordResetRepo.InvalidateUnusedForUser(ctx, user.ID); err != nil {
		return logger.Error(ctx, err, "Error while invalidating previous password resets")
	}

	now := time.Now().UTC()
	reset := &model.PasswordReset{
		PasswordResetID: uuid.New(),
		UserID:          user.ID,
		TokenHash:       tokenHash,
		ExpiresAt:       now.Add(passwordResetTTL),
		CreatedAt:       now,
	}

	if err = svc.passwordResetRepo.Insert(ctx, reset); err != nil {
		return logger.Error(ctx, err, "Error while inserting password reset")
	}

	dbPayload := map[string]any{
		"user_id":           user.ID,
		"password_reset_id": reset.PasswordResetID.String(),
		"raw_token":         rawToken,
		"ip":                ip,
		"user_agent":        userAgent,
	}
	queuePayload := map[string]any{
		"user_id":           user.ID,
		"password_reset_id": reset.PasswordResetID.String(),
		"ip":                ip,
		"user_agent":        userAgent,
	}

	if pubErr := svc.eventPublisher.PublishWithQueuePayload(ctx, types.EventTypePasswordResetRequested, dbPayload, queuePayload); pubErr != nil {
		_ = logger.Error(ctx, pubErr, "failed to publish user.password_reset_requested event")
	}

	return nil
}

func (svc *authService) ValidateResetToken(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidResetToken
	}
	tokenHash := hashPasswordResetToken(token)
	reset, err := svc.passwordResetRepo.GetValidByTokenHash(ctx, tokenHash)
	if err != nil {
		return logger.Error(ctx, err, "Error while validating password reset token")
	}
	if reset == nil {
		return ErrInvalidResetToken
	}
	return nil
}

func (svc *authService) ResetPassword(ctx context.Context, token, newPassword, ip, userAgent string) error {
	tokenHash := hashPasswordResetToken(token)
	reset, err := svc.passwordResetRepo.GetValidByTokenHash(ctx, tokenHash)
	if err != nil {
		return logger.Error(ctx, err, "Error while fetching password reset")
	}
	if reset == nil {
		return ErrInvalidResetToken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return logger.Error(ctx, err, "Error while hashing new password")
	}

	if err = svc.userRepo.UpdatePasswordHash(ctx, reset.UserID, string(hash)); err != nil {
		return logger.Error(ctx, err, "Error while updating password")
	}

	usedAt := time.Now().UTC()
	if err = svc.passwordResetRepo.MarkUsed(ctx, reset.PasswordResetID.String(), usedAt); err != nil {
		return logger.Error(ctx, err, "Error while marking password reset used")
	}

	if pubErr := svc.eventPublisher.Publish(ctx, types.EventTypePasswordChanged, map[string]any{
		"user_id":    reset.UserID,
		"ip":         ip,
		"user_agent": userAgent,
	}); pubErr != nil {
		_ = logger.Error(ctx, pubErr, "failed to publish user.password_changed event")
	}

	return nil
}

func generatePasswordResetToken() (rawToken, tokenHash string, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", err
	}
	rawToken = base64.RawURLEncoding.EncodeToString(buf)
	tokenHash = hashPasswordResetToken(rawToken)
	return rawToken, tokenHash, nil
}

func hashPasswordResetToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
