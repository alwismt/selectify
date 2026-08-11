package model_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"alwis.dev/selectify/internal/model"
)

func TestUserSession_IsExpired(t *testing.T) {
	now := time.Now().UTC()

	t.Run("valid", func(t *testing.T) {
		s := &model.LoggedInSession{
			ExpiresAt:         now.Add(time.Hour),
			AbsoluteExpiresAt: now.Add(24 * time.Hour),
		}
		require.False(t, s.IsExpired())
	})

	t.Run("idle expired", func(t *testing.T) {
		s := &model.LoggedInSession{
			ExpiresAt:         now.Add(-time.Minute),
			AbsoluteExpiresAt: now.Add(24 * time.Hour),
		}
		require.True(t, s.IsExpired())
	})

	t.Run("absolute expired", func(t *testing.T) {
		s := &model.LoggedInSession{
			ExpiresAt:         now.Add(time.Hour),
			AbsoluteExpiresAt: now.Add(-time.Minute),
		}
		require.True(t, s.IsExpired())
	})

	t.Run("revoked", func(t *testing.T) {
		revoked := now
		s := &model.LoggedInSession{
			ExpiresAt:         now.Add(time.Hour),
			AbsoluteExpiresAt: now.Add(24 * time.Hour),
			RevokedAt:         &revoked,
		}
		require.True(t, s.IsExpired())
	})
}
