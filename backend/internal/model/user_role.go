package model

import (
	"time"

	"alwis.dev/selectify/internal/types"
)

type UserRole struct {
	ID           int                           `db:"id" json:"-"`
	UserID       uint                          `db:"user_id" json:"-"`
	Role         types.UserRole                `db:"role" json:"role"`
	MerchantRole *types.MerchantRolePermission `db:"merchant_role" json:"merchant_role,omitempty"`
	MerchantID   *uint                         `db:"merchant_id" json:"merchant_id,omitempty"`
	CreatedAt    time.Time                     `db:"created_at" json:"-"`
}
