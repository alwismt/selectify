package repo_test

import (
	"fmt"
	"math/rand"
	"testing"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/types"

	"github.com/stretchr/testify/require"
)

func TestUserRepo_AddUserWithTx(t *testing.T) {
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
		require.NotZero(t, user.ID)
		require.Equal(t, user.Email, user.Email)
		require.Equal(t, user.FirstName, user.FirstName)
		require.Equal(t, user.LastName, user.LastName)
		require.Equal(t, user.Phone, user.Phone)
		require.Equal(t, types.UserStatusActive, user.Status)

		tx.CanCommit = true
		tx.End()

		t.Cleanup(func() {
			_, _ = dbConn.Exec("DELETE FROM users WHERE id = $1", user.ID)
		})
	})

	t.Run("With Duplicate Email", func(t *testing.T) {
		tx, err := transactionRepo.BeginTransaction(ctx)
		require.NoError(t, err)
		require.NotNil(t, tx)

		user := &model.User{
			Email:        fmt.Sprintf("duplicate%d@example.com", rand.Int()),
			PasswordHash: "hashed_password_123",
			FirstName:    "Test",
			LastName:     "User",
			Phone:        fmt.Sprintf("370123456%d", rand.Int()%10000),
			Status:       types.UserStatusActive,
		}

		err = userRepo.AddUserWithTx(ctx, tx.Transaction, user)
		require.NoError(t, err)
		require.NotZero(t, user.ID)

		tx.CanCommit = true
		tx.End()

		t.Cleanup(func() {
			_, _ = dbConn.Exec("DELETE FROM users WHERE id = $1", user.ID)
		})

		tx2, err := transactionRepo.BeginTransaction(ctx)
		require.NoError(t, err)
		require.NotNil(t, tx2)

		duplicateUser := &model.User{
			Email:        user.Email,
			PasswordHash: "hashed_password_456",
			FirstName:    "Duplicate",
			LastName:     "User",
			Phone:        fmt.Sprintf("370123456%d", rand.Int()%10000),
			Status:       types.UserStatusActive,
		}

		err = userRepo.AddUserWithTx(ctx, tx2.Transaction, duplicateUser)
		require.Error(t, err)
		require.Equal(t, httpx.ErrUserAlreadyExists, err)

		tx2.End()
	})
}

func TestUserRepo_GetUserByEmail(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		user, err := userRepo.GetUserByEmail(ctx, testUser.Email)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, testUser.ID, user.ID)
		require.Equal(t, testUser.Email, user.Email)
		require.Equal(t, testUser.FirstName, user.FirstName)
		require.Equal(t, testUser.LastName, user.LastName)
		require.Equal(t, testUser.Phone, user.Phone)
		require.Equal(t, testUser.Status, user.Status)
	})

	t.Run("NotFound", func(t *testing.T) {
		user, err := userRepo.GetUserByEmail(ctx, "nonexistent@example.com")

		require.Error(t, err)
		require.Equal(t, httpx.ErrUserNotFound, err)
		require.Nil(t, user)
	})

	t.Run("EmptyEmail", func(t *testing.T) {
		user, err := userRepo.GetUserByEmail(ctx, "")

		require.Error(t, err)
		require.Nil(t, user)
	})
}
