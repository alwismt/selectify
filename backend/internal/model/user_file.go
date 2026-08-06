package model

import (
	"time"

	"github.com/google/uuid"
)

type UserFile struct {
	ID          uuid.UUID  `db:"user_file_id" json:"id"`
	UserID      uint       `db:"user_id" json:"user_id"`
	ContentType string     `db:"content_type" json:"content_type"`
	CreatedAt   *time.Time `db:"created_at" json:"-"`
	UpdatedAt   *time.Time `db:"updated_at" json:"-"`
}
