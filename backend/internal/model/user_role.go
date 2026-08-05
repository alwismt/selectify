package model

import (
	"alwis.dev/selectify/internal/types"
	"database/sql"
	"time"
)

type UserRole struct {
	ID        int            `db:"id" json:"-"`
	UserID    uint           `db:"user_id" json:"user_id"`
	Role      types.UserRole `db:"role" json:"role"`
	ScopeType sql.NullString `db:"scope_type" json:"scope_type,omitempty"`
	ScopeID   int            `db:"scope_id" json:"scope_id,omitempty"`
	CreatedAt time.Time      `db:"created_at" json:"-"`
}
