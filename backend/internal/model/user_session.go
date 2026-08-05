package model

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	SessionId uuid.UUID  `db:"session_id" json:"-"`
	UserId    uint       `db:"user_id" json:"userId"`
	UserAgent string     `db:"user_agent" json:"userAgent"`
	IpAddress string     `db:"ip_address" json:"ipAddress"`
	IssuedAt  time.Time  `db:"issued_at" json:"-"`
	ExpiresAt time.Time  `db:"expires_at" json:"-"`
	RevokedAt *time.Time `db:"revoked_at" json:"-"`
	User      *User      `db:"-" json:"-"`
}

func (u *UserSession) IsExpired() bool {
	if u.RevokedAt != nil {
		return true
	}
	return time.Now().After(u.ExpiresAt)
}
