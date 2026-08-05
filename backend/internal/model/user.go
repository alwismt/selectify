package model

import (
	"time"

	"alwis.dev/selectify/internal/types"
)

type User struct {
	ID              uint             `db:"id" json:"id"`
	Email           string           `db:"email" json:"email"`
	FirstName       string           `db:"first_name" json:"first_name"`
	LastName        string           `db:"last_name" json:"last_name"`
	Phone           string           `db:"phone" json:"phone"`
	Status          types.UserStatus `db:"status" json:"status"`
	PasswordHash    string           `db:"password_hash" json:"-"`
	EmailVerifiedAt *time.Time       `db:"email_verified_at" json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time       `db:"last_login_at" json:"last_login_at,omitempty"`
	CreatedAt       *time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       *time.Time       `db:"updated_at" json:"-"`
	DeletedAt       *time.Time       `db:"deleted_at" json:"deleted_at,omitempty"`
}
