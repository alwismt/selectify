package repo_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"

	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/types"

	"github.com/stretchr/testify/require"
)

func TestUserSessionRepo_InsertUserSession(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tx, err := transactionRepo.BeginTransaction(ctx)
		require.NoError(t, err)
		require.NotNil(t, tx)

		user := &model.User{
			Email:        fmt.Sprintf("testuser%d@example.com", rand.Int()),
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

		t.Cleanup(func() {
			_, _ = dbConn.Exec("DELETE FROM user_session WHERE user_id = $1", user.ID)
			_, _ = dbConn.Exec("DELETE FROM users WHERE id = $1", user.ID)
		})

		session := &model.UserSession{
			UserId:    user.ID,
			SessionId: uuid.New(),
			ExpiresAt: time.Now().Add(24 * time.Hour),
			UserAgent: "Mozilla/5.0",
			IpAddress: "192.168.1.1",
		}

		err = userSessionRepo.InsertUserSession(ctx, session)
		require.NoError(t, err)

		var count int
		err = dbConn.Get(&count, "SELECT COUNT(*) FROM user_session WHERE user_id = $1 AND session_id = $2", user.ID, session.SessionId)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("With Existing User", func(t *testing.T) {
		session := &model.UserSession{
			UserId:    testUser.ID,
			SessionId: uuid.New(),
			ExpiresAt: time.Now().Add(24 * time.Hour),
			UserAgent: "Mozilla/5.0",
			IpAddress: "192.168.1.1",
		}

		err := userSessionRepo.InsertUserSession(ctx, session)
		require.NoError(t, err)

		t.Cleanup(func() {
			_, _ = dbConn.Exec("DELETE FROM user_session WHERE session_id = $1", session.SessionId)
		})

		var count int
		err = dbConn.Get(&count, "SELECT COUNT(*) FROM user_session WHERE user_id = $1 AND session_id = $2", testUser.ID, session.SessionId)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})
}
