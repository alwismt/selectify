package service

import (
	"context"
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

var (
	ErrUserDeactivated = errors.New("user is deactivated")
)

type AuthService interface {
	RegisterUser(ctx context.Context, req *request.UserRegisterRequest, userAgent, ip string) (*model.UserSession, error)
	LoginUser(ctx context.Context, req *request.LoginRequest, userAgent, ip string) (*model.UserSession, error)
	UserLogout(ctx context.Context, sesId uuid.UUID) (err error)
}

type authService struct {
	jwtService JWTService

	txRepo       repo.TransactionRepo
	userRepo     repo.UserRepo
	userRoleRepo repo.UserRoleRepo
	sessionRepo  repo.UserSessionRepo
}

func NewAuthService(jwtService JWTService, userRepo repo.UserRepo, txRepo repo.TransactionRepo, userRole repo.UserRoleRepo,
	sessionRepo repo.UserSessionRepo) AuthService {
	return &authService{
		jwtService: jwtService,

		txRepo:       txRepo,
		userRepo:     userRepo,
		userRoleRepo: userRole,
		sessionRepo:  sessionRepo,
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

	return svc.newUserSession(ctx, user.ID, userAgent, ip)
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
