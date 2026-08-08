package model

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	SessionId         uuid.UUID  `db:"session_id" json:"-"`
	UserId            uint       `db:"user_id" json:"userId"`
	UserAgent         string     `db:"user_agent" json:"userAgent"`
	IpAddress         string     `db:"ip_address" json:"ipAddress"`
	IssuedAt          time.Time  `db:"issued_at" json:"-"`
	ExpiresAt         time.Time  `db:"expires_at" json:"-"`
	AbsoluteExpiresAt time.Time  `db:"absolute_expires_at" json:"-"`
	LastActivityAt    time.Time  `db:"last_activity_at" json:"-"`
	RememberMe        bool       `db:"remember_me" json:"-"`
	SessionTokenHash  string     `db:"session_token_hash" json:"-"`
	DeviceId          *uuid.UUID `db:"device_id" json:"-"`
	RevokedAt         *time.Time `db:"revoked_at" json:"-"`
	User              *User      `db:"-" json:"-"`
	UserRole          *UserRole  `db:"-" json:"-"`

	// RawSessionToken is the bearer secret for the cookie; never persisted.
	RawSessionToken string `db:"-" json:"-"`
	// RawDeviceToken is set when a new device cookie should be issued; never persisted.
	RawDeviceToken string `db:"-" json:"-"`
}

func (u *UserSession) IsExpired() bool {
	if u.RevokedAt != nil {
		return true
	}
	now := time.Now().UTC()
	if now.After(u.ExpiresAt) {
		return true
	}
	if now.After(u.AbsoluteExpiresAt) {
		return true
	}
	return false
}
