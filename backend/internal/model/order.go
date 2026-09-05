package model

import (
	"encoding/json"
	"time"

	"alwis.dev/selectify/internal/types"
)

type Orders []Order

type Order struct {
	ID           uint              `db:"id" json:"id"`
	UserID       uint              `db:"user_id" json:"user_id"`
	Status       types.OrderStatus `db:"status" json:"status"`
	Currency     string            `db:"currency" json:"currency"`
	Subtotal     uint64            `db:"subtotal" json:"subtotal"`
	Shipping     uint64            `db:"shipping" json:"shipping"`
	Discount     uint64            `db:"discount" json:"discount"`
	Total        uint64            `db:"total" json:"total"`
	CreatedAt    time.Time         `db:"created_at" json:"-"`
	UpdatedAt    time.Time         `db:"updated_at" json:"-"`
	Items        OrderItems        `db:"-" json:"items,omitempty"`
	ClientSecret string            `db:"-" json:"client_secret,omitempty"`
}

type OrderItems []OrderItem

type OrderItem struct {
	ID         uint            `db:"id" json:"id"`
	OrderID    uint            `db:"order_id" json:"order_id"`
	VariantID  uint            `db:"variant_id" json:"variant_id"`
	SKU        string          `db:"sku" json:"sku"`
	UnitPrice  uint64          `db:"unit_price" json:"unit_price"`
	Quantity   uint            `db:"quantity" json:"quantity"`
	Attributes json.RawMessage `db:"attributes" json:"attributes"` // JSONB stored as string
	CreatedAt  time.Time       `db:"created_at" json:"created_at"`
}
