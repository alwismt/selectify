package model

import (
	"time"
)

type Cart struct {
	ID        uint       `db:"id" json:"id"`
	UserID    uint       `db:"user_id" json:"user_id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	Items     []CartItem `db:"-" json:"items,omitempty"`
}

type CartItems []CartItem

type CartItem struct {
	ID        uint      `db:"id" json:"id"`
	CartID    uint      `db:"cart_id" json:"cart_id"`
	VariantID uint      `db:"variant_id" json:"variant_id"`
	Quantity  uint      `db:"quantity" json:"quantity"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Cart      *Cart     `db:"-" json:"-"`
}

func (c CartItems) GetVariantIDs() []uint {
	var ids []uint

	if len(c) == 0 {
		return ids
	}

	for _, ci := range c {
		ids = append(ids, ci.VariantID)
	}
	return ids
}
