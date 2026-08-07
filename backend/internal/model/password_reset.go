package model

import (
	"time"

	"github.com/google/uuid"
)

type PasswordReset struct {
	PasswordResetID uuid.UUID  `db:"password_reset_id" json:"-"`
	UserID          uint       `db:"user_id" json:"-"`
	TokenHash       string     `db:"token_hash" json:"-"`
	ExpiresAt       time.Time  `db:"expires_at" json:"-"`
	UsedAt          *time.Time `db:"used_at" json:"-"`
	CreatedAt       time.Time  `db:"created_at" json:"-"`
}
