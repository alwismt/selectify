package service

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"

	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
)

type UserService interface {
	GetUserImage(ctx context.Context, user *model.User) (*model.UserFile, error)
	UpsertUserImage(ctx context.Context, user *model.User, file io.Reader, contentType string) (*model.UserFile, error)
	DeleteUserImage(ctx context.Context, user *model.User) error
}

type userService struct {
	storageService StorageService

	txRepo       repo.TransactionRepo
	userFileRepo repo.UserFileRepo
}

func NewUserService(storageService StorageService, txRepo repo.TransactionRepo, userFileRepo repo.UserFileRepo) UserService {
	return &userService{
		storageService: storageService,

		txRepo:       txRepo,
		userFileRepo: userFileRepo,
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

func userImageObjectKey(userFileID uuid.UUID) string {
	return fmt.Sprintf("users/%s", userFileID.String())
}
