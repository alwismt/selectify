package service

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

type mockStorageService struct {
	gotObjectKey   string
	gotContentType string
	gotBody        []byte
	gotDeletedKey  string
	uploadErr      error
	deleteErr      error
}

func (m *mockStorageService) UploadFile(_ context.Context, file io.Reader, objectKey string, contentType string) error {
	m.gotObjectKey = objectKey
	m.gotContentType = contentType

	body, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	m.gotBody = body

	return m.uploadErr
}

func (m *mockStorageService) DeleteFile(_ context.Context, objectKey string) error {
	m.gotDeletedKey = objectKey
	return m.deleteErr
}

type mockTransactionRepo struct {
	tx *db.Transaction
}

func (m *mockTransactionRepo) BeginTransaction(_ context.Context) (*db.Transaction, error) {
	if m.tx == nil {
		m.tx = &db.Transaction{}
	}
	return m.tx, nil
}

type mockUserFileRepo struct {
	gotFile         *model.UserFile
	getByUserIDFile *model.UserFile
	deleteByUserID  uint
}

func (m *mockUserFileRepo) AddUserFileWithTx(_ context.Context, _ sqlx.ExecerContext, file *model.UserFile) error {
	fileCopy := *file
	m.gotFile = &fileCopy
	return nil
}

func (m *mockUserFileRepo) GetByUserID(_ context.Context, _ uint) (*model.UserFile, error) {
	return m.getByUserIDFile, nil
}

func (m *mockUserFileRepo) DeleteByUserIDWithTx(_ context.Context, _ sqlx.ExecerContext, userID uint) error {
	m.deleteByUserID = userID
	return nil
}

func TestUserService_GetUserImage(t *testing.T) {
	existing := &model.UserFile{ID: uuid.New(), UserID: 1, ContentType: "image/png"}
	userFileRepo := &mockUserFileRepo{getByUserIDFile: existing}
	svc := NewUserService(&mockStorageService{}, &mockTransactionRepo{}, userFileRepo, nil)

	userFile, err := svc.GetUserImage(context.Background(), &model.User{ID: 1})
	require.NoError(t, err)
	require.Equal(t, existing.ID, userFile.ID)
}

func TestUserService_GetUserImage_NoExistingFile(t *testing.T) {
	svc := NewUserService(&mockStorageService{}, &mockTransactionRepo{}, &mockUserFileRepo{}, nil)

	userFile, err := svc.GetUserImage(context.Background(), &model.User{ID: 1})
	require.NoError(t, err)
	require.Nil(t, userFile)
}

func TestUserService_UpsertUserImage(t *testing.T) {
	storage := &mockStorageService{}
	txRepo := &mockTransactionRepo{}
	userFileRepo := &mockUserFileRepo{}
	svc := NewUserService(storage, txRepo, userFileRepo, nil)
	user := &model.User{ID: 1}

	userFile, err := svc.UpsertUserImage(context.Background(), user, strings.NewReader("image-bytes"), "image/png")
	require.NoError(t, err)
	require.NotNil(t, userFile)
	require.Equal(t, uint(1), userFile.UserID)
	require.Equal(t, "image/png", userFile.ContentType)
	require.Equal(t, "image/png", storage.gotContentType)
	require.Equal(t, "image-bytes", string(storage.gotBody))
	require.Equal(t, "users/"+userFile.ID.String(), storage.gotObjectKey)
	require.Empty(t, storage.gotDeletedKey)
	require.NotNil(t, userFileRepo.gotFile)
	require.Equal(t, userFile.ID, userFileRepo.gotFile.ID)
	require.True(t, txRepo.tx.CanCommit)
}

func TestUserService_UpsertUserImage_RollsBackWhenUploadFails(t *testing.T) {
	storage := &mockStorageService{uploadErr: errors.New("upload failed")}
	txRepo := &mockTransactionRepo{}
	oldFile := &model.UserFile{ID: uuid.New(), UserID: 1, ContentType: "image/jpeg"}
	userFileRepo := &mockUserFileRepo{getByUserIDFile: oldFile}
	svc := NewUserService(storage, txRepo, userFileRepo, nil)
	user := &model.User{ID: 1}

	userFile, err := svc.UpsertUserImage(context.Background(), user, strings.NewReader("image-bytes"), "image/png")
	require.Error(t, err)
	require.Nil(t, userFile)
	require.NotNil(t, userFileRepo.gotFile)
	require.Empty(t, storage.gotDeletedKey)
	require.False(t, txRepo.tx.CanCommit)
}

func TestUserService_UpsertUserImage_ReplacesExistingImage(t *testing.T) {
	oldFile := &model.UserFile{ID: uuid.New(), UserID: 1, ContentType: "image/jpeg"}
	storage := &mockStorageService{}
	txRepo := &mockTransactionRepo{}
	userFileRepo := &mockUserFileRepo{getByUserIDFile: oldFile}
	svc := NewUserService(storage, txRepo, userFileRepo, nil)
	user := &model.User{ID: 1}

	userFile, err := svc.UpsertUserImage(context.Background(), user, strings.NewReader("new-image-bytes"), "image/png")
	require.NoError(t, err)
	require.NotNil(t, userFile)
	require.NotEqual(t, oldFile.ID, userFile.ID)
	require.Equal(t, "users/"+userFile.ID.String(), storage.gotObjectKey)
	require.Equal(t, "users/"+oldFile.ID.String(), storage.gotDeletedKey)
	require.True(t, txRepo.tx.CanCommit)
}

func TestUserService_DeleteUserImage(t *testing.T) {
	userFile := &model.UserFile{ID: uuid.New(), UserID: 1, ContentType: "image/png"}
	storage := &mockStorageService{}
	txRepo := &mockTransactionRepo{}
	userFileRepo := &mockUserFileRepo{getByUserIDFile: userFile}
	svc := NewUserService(storage, txRepo, userFileRepo, nil)

	err := svc.DeleteUserImage(context.Background(), &model.User{ID: 1})
	require.NoError(t, err)
	require.Equal(t, "users/"+userFile.ID.String(), storage.gotDeletedKey)
	require.Equal(t, uint(1), userFileRepo.deleteByUserID)
	require.True(t, txRepo.tx.CanCommit)
}

func TestUserService_DeleteUserImage_NoExistingFile(t *testing.T) {
	storage := &mockStorageService{}
	txRepo := &mockTransactionRepo{}
	userFileRepo := &mockUserFileRepo{}
	svc := NewUserService(storage, txRepo, userFileRepo, nil)

	err := svc.DeleteUserImage(context.Background(), &model.User{ID: 1})
	require.NoError(t, err)
	require.Empty(t, storage.gotDeletedKey)
	require.Zero(t, userFileRepo.deleteByUserID)
	require.Nil(t, txRepo.tx)
}

func TestUserService_DeleteUserImage_RollsBackWhenStorageDeleteFails(t *testing.T) {
	userFile := &model.UserFile{ID: uuid.New(), UserID: 1, ContentType: "image/png"}
	storage := &mockStorageService{deleteErr: errors.New("delete failed")}
	txRepo := &mockTransactionRepo{}
	userFileRepo := &mockUserFileRepo{getByUserIDFile: userFile}
	svc := NewUserService(storage, txRepo, userFileRepo, nil)

	err := svc.DeleteUserImage(context.Background(), &model.User{ID: 1})
	require.Error(t, err)
	require.Equal(t, "users/"+userFile.ID.String(), storage.gotDeletedKey)
	require.Zero(t, userFileRepo.deleteByUserID)
	require.False(t, txRepo.tx.CanCommit)
}
