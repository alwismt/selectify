package model

import (
	"time"

	"github.com/google/uuid"
)

type UserDevice struct {
	DeviceId        uuid.UUID `db:"device_id" json:"-"`
	UserId          uint      `db:"user_id" json:"userId"`
	DeviceTokenHash string    `db:"device_token_hash" json:"-"`
	UserAgent       string    `db:"user_agent" json:"userAgent"`
	LastIP          string    `db:"last_ip" json:"lastIp"`
	FirstSeenAt     time.Time `db:"first_seen_at" json:"-"`
	LastSeenAt      time.Time `db:"last_seen_at" json:"-"`

	// RawDeviceToken is the bearer secret for the device cookie; never persisted.
	RawDeviceToken string `db:"-" json:"-"`
}
