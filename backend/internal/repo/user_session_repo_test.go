package repo_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/security"
	"alwis.dev/selectify/internal/service"
	"alwis.dev/selectify/internal/types"
)

func TestUserSessionRepo_InsertAndGetByTokenHash(t *testing.T) {
	tx, err := transactionRepo.BeginTransaction(ctx)
	require.NoError(t, err)

	user := &model.User{
		Email:        fmt.Sprintf("sessionuser%d@example.com", rand.Int()),
		PasswordHash: "hashed_password_123",
		FirstName:    "Test",
		LastName:     "User",
		Phone:        fmt.Sprintf("370123456%d", rand.Int()%10000),
		Status:       types.UserStatusActive,
	}
	err = userRepo.AddUserWithTx(ctx, tx.Transaction, user)
	require.NoError(t, err)
	tx.CanCommit = true
	tx.End()

	raw, hash, err := security.NewToken()
	require.NoError(t, err)

	now := time.Now().UTC()
	session := &model.LoggedInSession{
		UserId:            user.ID,
		SessionId:         uuid.New(),
		SessionTokenHash:  hash,
		ExpiresAt:         now.Add(6 * time.Hour),
		AbsoluteExpiresAt: now.Add(7 * 24 * time.Hour),
		LastActivityAt:    now.Add(-2 * time.Hour),
		RememberMe:        false,
		UserAgent:         "Mozilla/5.0",
		IpAddress:         "192.168.1.1",
	}

	err = userSessionRepo.InsertUserSession(ctx, session)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = dbConn.Exec("DELETE FROM user_session WHERE user_id = $1", user.ID)
		_, _ = dbConn.Exec("DELETE FROM users WHERE id = $1", user.ID)
	})

	got, err := userSessionRepo.GetByTokenHash(ctx, security.HashToken(raw))
	require.NoError(t, err)
	require.Equal(t, session.SessionId, got.SessionId)
	require.Equal(t, hash, got.SessionTokenHash)
	require.False(t, got.RememberMe)

	// Raw secret must not be findable via session_id string
	_, err = userSessionRepo.GetByTokenHash(ctx, session.SessionId.String())
	require.Error(t, err)
}

func TestUserSessionRepo_RenewIfStale(t *testing.T) {
	raw, hash, err := security.NewToken()
	require.NoError(t, err)
	_ = raw

	now := time.Now().UTC()
	session := &model.LoggedInSession{
		UserId:            testUser.ID,
		SessionId:         uuid.New(),
		SessionTokenHash:  hash,
		ExpiresAt:         now.Add(1 * time.Hour),
		AbsoluteExpiresAt: now.Add(7 * 24 * time.Hour),
		LastActivityAt:    now.Add(-2 * time.Hour),
		RememberMe:        false,
		UserAgent:         "Mozilla/5.0",
		IpAddress:         "10.0.0.1",
	}
	require.NoError(t, userSessionRepo.InsertUserSession(ctx, session))
	t.Cleanup(func() {
		_, _ = dbConn.Exec("DELETE FROM user_session WHERE session_id = $1", session.SessionId)
	})

	renewed, err := userSessionRepo.RenewIfStale(ctx, session.SessionId, service.SessionIdleTTL, service.ActivityThrottle)
	require.NoError(t, err)
	require.NotNil(t, renewed)
	require.True(t, renewed.ExpiresAt.After(now.Add(5*time.Hour)))
	require.True(t, renewed.LastActivityAt.After(now.Add(-time.Minute)))

	// Within throttle window — no second write
	again, err := userSessionRepo.RenewIfStale(ctx, session.SessionId, service.SessionIdleTTL, service.ActivityThrottle)
	require.NoError(t, err)
	require.Nil(t, again)
}

func TestUserSessionRepo_RememberMeNotRenewed(t *testing.T) {
	_, hash, err := security.NewToken()
	require.NoError(t, err)

	now := time.Now().UTC()
	session := &model.LoggedInSession{
		UserId:            testUser.ID,
		SessionId:         uuid.New(),
		SessionTokenHash:  hash,
		ExpiresAt:         now.Add(30 * 24 * time.Hour),
		AbsoluteExpiresAt: now.Add(30 * 24 * time.Hour),
		LastActivityAt:    now.Add(-2 * time.Hour),
		RememberMe:        true,
		UserAgent:         "Mozilla/5.0",
		IpAddress:         "10.0.0.2",
	}
	require.NoError(t, userSessionRepo.InsertUserSession(ctx, session))
	t.Cleanup(func() {
		_, _ = dbConn.Exec("DELETE FROM user_session WHERE session_id = $1", session.SessionId)
	})

	renewed, err := userSessionRepo.RenewIfStale(ctx, session.SessionId, service.SessionIdleTTL, service.ActivityThrottle)
	require.NoError(t, err)
	require.Nil(t, renewed)
}

func TestUserSessionRepo_RevokeSession(t *testing.T) {
	_, hash, err := security.NewToken()
	require.NoError(t, err)

	now := time.Now().UTC()
	session := &model.LoggedInSession{
		UserId:            testUser.ID,
		SessionId:         uuid.New(),
		SessionTokenHash:  hash,
		ExpiresAt:         now.Add(6 * time.Hour),
		AbsoluteExpiresAt: now.Add(7 * 24 * time.Hour),
		LastActivityAt:    now,
		UserAgent:         "Mozilla/5.0",
		IpAddress:         "10.0.0.3",
	}
	require.NoError(t, userSessionRepo.InsertUserSession(ctx, session))
	t.Cleanup(func() {
		_, _ = dbConn.Exec("DELETE FROM user_session WHERE session_id = $1", session.SessionId)
	})

	require.NoError(t, userSessionRepo.RevokeSession(ctx, session.SessionId))
	_, err = userSessionRepo.GetByTokenHash(ctx, hash)
	require.Error(t, err)
}

func TestUserDeviceRepo_InsertAndReuse(t *testing.T) {
	raw, hash, err := security.NewToken()
	require.NoError(t, err)

	now := time.Now().UTC()
	device := &model.UserDevice{
		DeviceId:        uuid.New(),
		UserId:          testUser.ID,
		DeviceTokenHash: hash,
		UserAgent:       "Mozilla/5.0",
		LastIP:          "127.0.0.1",
		FirstSeenAt:     now,
		LastSeenAt:      now.Add(-2 * time.Hour),
	}
	require.NoError(t, userDeviceRepo.Insert(ctx, device))
	t.Cleanup(func() {
		_, _ = dbConn.Exec("DELETE FROM user_device WHERE device_id = $1", device.DeviceId)
	})

	got, err := userDeviceRepo.GetByTokenHash(ctx, security.HashToken(raw))
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, device.DeviceId, got.DeviceId)

	require.NoError(t, userDeviceRepo.TouchIfStale(ctx, device.DeviceId, "UA2", "1.2.3.4", service.ActivityThrottle))
	got2, err := userDeviceRepo.GetByTokenHash(ctx, hash)
	require.NoError(t, err)
	require.Equal(t, "UA2", got2.UserAgent)
	require.Equal(t, "1.2.3.4", got2.LastIP)
}
