package repo_test

import (
	"fmt"
	"math/rand"
	"testing"

	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/types"

	"github.com/stretchr/testify/require"
)

func TestUserRoleRepo_InsertUserRoleForCustomerWithTx(t *testing.T) {
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

		role := &model.UserRole{
			UserID: user.ID,
			Role:   types.RoleCustomer,
		}

		err = userRoleRepo.InsertUserRoleForCustomerWithTx(ctx, tx.Transaction, role)
		require.NoError(t, err)

		tx.CanCommit = true
		tx.End()

		t.Cleanup(func() {
			_, _ = dbConn.Exec("DELETE FROM user_role WHERE user_id = $1", user.ID)
			_, _ = dbConn.Exec("DELETE FROM users WHERE id = $1", user.ID)
		})

		var count int
		err = dbConn.Get(&count, "SELECT COUNT(*) FROM user_role WHERE user_id = $1 AND role = $2", user.ID, types.RoleCustomer)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})
}
