package model

import (
	"time"

	"alwis.dev/selectify/internal/types"
)

type CheckoutSession struct {
	ID        uint                  `db:"id" json:"id"`
	OrderID   uint                  `db:"order_id" json:"order_id,omitempty"`
	Status    types.PaymentStatus   `db:"status" json:"status"`
	ExpiresAt time.Time             `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time             `db:"created_at" json:"created_at"`
	UpdatedAt time.Time             `db:"updated_at" json:"updated_at"`
	Items     []CheckoutSessionItem `db:"-" json:"items,omitempty"`
}

type CheckoutSessionItems []CheckoutSessionItem

type CheckoutSessionItem struct {
	CheckoutSessionID uint `db:"checkout_session_id" json:"checkout_session_id"`
	CartItemID        uint `db:"cart_item_id" json:"cart_item_id"`
}
