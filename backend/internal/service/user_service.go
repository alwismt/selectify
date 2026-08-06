package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	"alwis.dev/selectify/internal/email"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
	"alwis.dev/selectify/internal/types"
)

type UserService interface {
	GetUserImage(ctx context.Context, user *model.User) (*model.UserFile, error)
	UpsertUserImage(ctx context.Context, user *model.User, file io.Reader, contentType string) (*model.UserFile, error)
	DeleteUserImage(ctx context.Context, user *model.User) error
	ProcessUserLoggedIn(ctx context.Context, event *model.Event) error
}

type loginNotificationData struct {
	FirstName   string
	LastName    string
	LoggedInAt  string
	IP          string
	UserAgent   string
	Country     string
	City        string
	Subdivision string
	Timezone    string
}

type userService struct {
	storageService StorageService
	emailService   EmailService
	geoIPService   GeoIPService

	txRepo       repo.TransactionRepo
	userFileRepo repo.UserFileRepo
	userRepo     repo.UserRepo
}

func NewUserService(
	storageService StorageService,
	txRepo repo.TransactionRepo,
	userFileRepo repo.UserFileRepo,
	userRepo repo.UserRepo,
	emailService EmailService,
	geoIPService GeoIPService,
) UserService {
	return &userService{
		storageService: storageService,
		emailService:   emailService,
		geoIPService:   geoIPService,

		txRepo:       txRepo,
		userFileRepo: userFileRepo,
		userRepo:     userRepo,
	}
}

func (s *userService) GetUserImage(ctx context.Context, user *model.User) (*model.UserFile, error) {
	if user == nil {
		return nil, fmt.Errorf("user is required")
	}

	userFile, err := s.userFileRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, logger.Error(ctx, err, "error get user file")
	}

	return userFile, nil
}

func (s *userService) UpsertUserImage(ctx context.Context, user *model.User, file io.Reader, contentType string) (*model.UserFile, error) {
	if user == nil {
		return nil, fmt.Errorf("user is required")
	}
	if file == nil {
		return nil, fmt.Errorf("file is required")
	}
	if contentType == "" {
		return nil, fmt.Errorf("content type is required")
	}

	existingFile, err := s.userFileRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, logger.Error(ctx, err, "error get user file")
	}

	var oldFileID *uuid.UUID
	if existingFile != nil {
		oldID := existingFile.ID
		oldFileID = &oldID
	}

	userFile := &model.UserFile{
		ID:          uuid.New(),
		UserID:      user.ID,
		ContentType: contentType,
		CreatedAt:   nil,
		UpdatedAt:   nil,
	}
	tx, err := s.txRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, logger.Error(ctx, err, "error begin transaction")
	}
	defer tx.End()

	if err = s.userFileRepo.AddUserFileWithTx(ctx, tx.Transaction, userFile); err != nil {
		return nil, logger.Error(ctx, err, "error add user file")
	}

	objectKey := userImageObjectKey(userFile.ID)
	if err = s.storageService.UploadFile(ctx, file, objectKey, contentType); err != nil {
		return nil, logger.Error(ctx, err, "error upload user file")
	}

	if oldFileID != nil {
		if err = s.storageService.DeleteFile(ctx, userImageObjectKey(*oldFileID)); err != nil {
			return nil, logger.Error(ctx, err, "error delete old user file from storage")
		}
	}

	tx.CanCommit = true
	return userFile, nil
}

func (s *userService) DeleteUserImage(ctx context.Context, user *model.User) error {
	if user == nil {
		return fmt.Errorf("user is required")
	}

	userFile, err := s.userFileRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		return logger.Error(ctx, err, "error get user file")
	}
	if userFile == nil {
		return nil
	}

	tx, err := s.txRepo.BeginTransaction(ctx)
	if err != nil {
		return logger.Error(ctx, err, "error begin transaction")
	}
	defer tx.End()

	if err = s.storageService.DeleteFile(ctx, userImageObjectKey(userFile.ID)); err != nil {
		return logger.Error(ctx, err, "error delete user file from storage")
	}
	if err = s.userFileRepo.DeleteByUserIDWithTx(ctx, tx.Transaction, user.ID); err != nil {
		return logger.Error(ctx, err, "error delete user file")
	}

	tx.CanCommit = true
	return nil
}

func (s *userService) ProcessUserLoggedIn(ctx context.Context, event *model.Event) error {
	if event == nil || event.Data == nil {
		return fmt.Errorf("event is required")
	}

	var payload struct {
		UserID    uint   `json:"user_id"`
		SessionID string `json:"session_id"`
		IP        string `json:"ip"`
		UserAgent string `json:"user_agent"`
	}
	if err := json.Unmarshal(event.Data.Payload, &payload); err != nil {
		return logger.Error(ctx, err, "failed to unmarshal user.logged_in payload")
	}
	if payload.UserID == 0 {
		return fmt.Errorf("user_id is required in user.logged_in payload")
	}

	user, err := s.userRepo.GetUserById(ctx, payload.UserID)
	if err != nil {
		return logger.Errorf(ctx, err, "failed to get user %d for event %s", payload.UserID, event.ID)
	}

	if user.Status != types.UserStatusActive {
		return fmt.Errorf("user %d is not active (status=%s)", user.ID, user.Status)
	}
	if user.Email == "" {
		return fmt.Errorf("user %d has no email address", user.ID)
	}
	if s.emailService == nil {
		return fmt.Errorf("email service is not configured")
	}

	ip := payload.IP
	if ip == "" {
		ip = "unknown"
	}
	userAgent := payload.UserAgent
	if userAgent == "" {
		userAgent = "unknown"
	}

	firstName := user.FirstName
	if firstName == "" {
		firstName = "there"
	}

	lastname := user.LastName
	if lastname == "" {
		lastname = ""
	}

	data := loginNotificationData{
		FirstName:  firstName,
		LastName:   lastname,
		LoggedInAt: time.Now().UTC().Format("2006-01-02 15:04:05 UTC"),
		IP:         ip,
		UserAgent:  userAgent,
	}
	if s.geoIPService != nil && ip != "unknown" {
		loc := s.geoIPService.Lookup(ip)
		data.Country = loc.Country
		data.City = loc.City
		data.Subdivision = loc.Subdivision
		data.Timezone = loc.Timezone
	}

	err = s.emailService.SendTemplate(ctx, TemplateMessage{
		To:       user.Email,
		Subject:  "New login to your Selectify account",
		Template: email.TemplateLoginNotification,
		Data:     data,
	})
	if err != nil {
		return logger.Errorf(ctx, err, "failed to send login notification email to user %d", user.ID)
	}

	return nil
}

func userImageObjectKey(userFileID uuid.UUID) string {
	return fmt.Sprintf("users/%s", userFileID.String())
}
