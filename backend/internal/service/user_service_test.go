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
	svc := NewUserService(&mockStorageService{}, &mockTransactionRepo{}, userFileRepo, nil, nil, nil)

	userFile, err := svc.GetUserImage(context.Background(), &model.User{ID: 1})
	require.NoError(t, err)
	require.Equal(t, existing.ID, userFile.ID)
}

func TestUserService_GetUserImage_NoExistingFile(t *testing.T) {
	svc := NewUserService(&mockStorageService{}, &mockTransactionRepo{}, &mockUserFileRepo{}, nil, nil, nil)

	userFile, err := svc.GetUserImage(context.Background(), &model.User{ID: 1})
	require.NoError(t, err)
	require.Nil(t, userFile)
}

func TestUserService_UpsertUserImage(t *testing.T) {
	storage := &mockStorageService{}
	txRepo := &mockTransactionRepo{}
	userFileRepo := &mockUserFileRepo{}
	svc := NewUserService(storage, txRepo, userFileRepo, nil, nil, nil)
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
	svc := NewUserService(storage, txRepo, userFileRepo, nil, nil, nil)
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
	svc := NewUserService(storage, txRepo, userFileRepo, nil, nil, nil)
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
	svc := NewUserService(storage, txRepo, userFileRepo, nil, nil, nil)

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
	svc := NewUserService(storage, txRepo, userFileRepo, nil, nil, nil)

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
	svc := NewUserService(storage, txRepo, userFileRepo, nil, nil, nil)

	err := svc.DeleteUserImage(context.Background(), &model.User{ID: 1})
	require.Error(t, err)
	require.Equal(t, "users/"+userFile.ID.String(), storage.gotDeletedKey)
	require.Zero(t, userFileRepo.deleteByUserID)
	require.False(t, txRepo.tx.CanCommit)
}

type mockEmailService struct {
	gotMsg    TemplateMessage
	sendErr   error
	sendCalls int
}

func (m *mockEmailService) Send(_ context.Context, _ EmailMessage) error {
	return m.sendErr
}

func (m *mockEmailService) SendTemplate(_ context.Context, msg TemplateMessage) error {
	m.sendCalls++
	m.gotMsg = msg
	return m.sendErr
}

type mockUserRepo struct {
	user *model.User
	err  error
}

func (m *mockUserRepo) AddUserWithTx(_ context.Context, _ sqlx.QueryerContext, _ *model.User) error {
	return nil
}

func (m *mockUserRepo) GetUserByEmail(_ context.Context, _ string) (*model.User, error) {
	return m.user, m.err
}

func (m *mockUserRepo) GetUserById(_ context.Context, _ uint) (*model.User, error) {
	return m.user, m.err
}

func (m *mockUserRepo) UpdatePasswordHash(_ context.Context, _ uint, _ string) error {
	return nil
}

func TestUserService_ProcessUserLoggedIn_SendsEmail(t *testing.T) {
	emailSvc := &mockEmailService{}
	userRepo := &mockUserRepo{
		user: &model.User{
			ID:        42,
			Email:     "ada@example.com",
			FirstName: "Ada",
			Status:    "active",
		},
	}
	svc := NewUserService(nil, nil, nil, userRepo, emailSvc, nil)

	event := &model.Event{
		ID: "evt-1",
		Data: &model.EventData{
			Payload: []byte(`{"user_id":42,"session_id":"sess-1","ip":"203.0.113.10","user_agent":"Mozilla/5.0"}`),
		},
	}

	err := svc.ProcessUserLoggedIn(context.Background(), event)
	require.NoError(t, err)
	require.Equal(t, 1, emailSvc.sendCalls)
	require.Equal(t, "ada@example.com", emailSvc.gotMsg.To)
	require.Equal(t, "New login to your Selectify account", emailSvc.gotMsg.Subject)
	require.Equal(t, "login_notification", emailSvc.gotMsg.Template)

	data, ok := emailSvc.gotMsg.Data.(loginNotificationData)
	require.True(t, ok)
	require.Equal(t, "Ada", data.FirstName)
	require.Equal(t, "203.0.113.10", data.IP)
	require.Equal(t, "Mozilla/5.0", data.UserAgent)
	require.NotEmpty(t, data.LoggedInAt)
	require.Equal(t, "Location unknown", data.Location)
}

type mockGeoIPService struct {
	gotIP string
	loc   Location
}

func (m *mockGeoIPService) Lookup(ip string) Location {
	m.gotIP = ip
	return m.loc
}

func (m *mockGeoIPService) Close() error {
	return nil
}

func TestUserService_ProcessUserLoggedIn_IncludesLocation(t *testing.T) {
	emailSvc := &mockEmailService{}
	geoIP := &mockGeoIPService{
		loc: Location{
			Country:     "United States",
			City:        "Mountain View",
			Subdivision: "California",
			Timezone:    "America/Los_Angeles",
		},
	}
	userRepo := &mockUserRepo{
		user: &model.User{
			ID:     42,
			Email:  "ada@example.com",
			Status: "active",
		},
	}
	svc := NewUserService(nil, nil, nil, userRepo, emailSvc, geoIP)

	event := &model.Event{
		ID: "evt-1",
		Data: &model.EventData{
			Payload: []byte(`{"user_id":42,"ip":"8.8.8.8","user_agent":"test"}`),
		},
	}

	err := svc.ProcessUserLoggedIn(context.Background(), event)
	require.NoError(t, err)
	require.Equal(t, "8.8.8.8", geoIP.gotIP)

	data, ok := emailSvc.gotMsg.Data.(loginNotificationData)
	require.True(t, ok)
	require.Equal(t, "Mountain View, United States", data.Location)
}

func TestUserService_ProcessUserLoggedIn_PropagatesSendError(t *testing.T) {
	emailSvc := &mockEmailService{sendErr: errors.New("smtp down")}
	userRepo := &mockUserRepo{
		user: &model.User{
			ID:     42,
			Email:  "ada@example.com",
			Status: "active",
		},
	}
	svc := NewUserService(nil, nil, nil, userRepo, emailSvc, nil)

	event := &model.Event{
		ID: "evt-1",
		Data: &model.EventData{
			Payload: []byte(`{"user_id":42,"ip":"1.2.3.4","user_agent":"test"}`),
		},
	}

	err := svc.ProcessUserLoggedIn(context.Background(), event)
	require.Error(t, err)
	require.Contains(t, err.Error(), "smtp down")
	require.Equal(t, 1, emailSvc.sendCalls)
}

func TestUserService_ProcessUserLoggedIn_RequiresEmailService(t *testing.T) {
	userRepo := &mockUserRepo{
		user: &model.User{
			ID:     42,
			Email:  "ada@example.com",
			Status: "active",
		},
	}
	svc := NewUserService(nil, nil, nil, userRepo, nil, nil)

	event := &model.Event{
		ID: "evt-1",
		Data: &model.EventData{
			Payload: []byte(`{"user_id":42}`),
		},
	}

	err := svc.ProcessUserLoggedIn(context.Background(), event)
	require.Error(t, err)
	require.Contains(t, err.Error(), "email service is not configured")
}
