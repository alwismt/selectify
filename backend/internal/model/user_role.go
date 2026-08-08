package model

import (
	"time"

	"alwis.dev/selectify/internal/types"
)

type UserRole struct {
	ID        int            `db:"id" json:"-"`
	UserID    uint           `db:"user_id" json:"-"`
	Role      types.UserRole `db:"role" json:"role"`
	ScopeType *string        `db:"scope_type" json:"scope_type,omitempty"`
	ScopeID   *uint          `db:"scope_id" json:"scope_id,omitempty"`
	CreatedAt time.Time      `db:"created_at" json:"-"`
}
