package model

import (
	"time"
)

type UserAddresses []UserAddress

type UserAddress struct {
	ID          uint      `db:"id" json:"id"`
	UserID      uint      `db:"user_id" json:"user_id"`
	Label       *string   `db:"label" json:"label,omitempty"`
	Phone       *string   `db:"phone" json:"phone,omitempty"`
	Line1       string    `db:"line1" json:"line1"`
	Line2       *string   `db:"line2" json:"line2,omitempty"`
	City        string    `db:"city" json:"city"`
	Region      *string   `db:"region" json:"region,omitempty"`
	PostalCode  string    `db:"postal_code" json:"postal_code"`
	CountryCode string    `db:"country_code" json:"country_code"`
	IsDefault   bool      `db:"is_default" json:"is_default"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type OrderAddresses []OrderAddress

type OrderAddress struct {
	ID          uint      `db:"id" json:"id"`
	OrderID     uint      `db:"order_id" json:"order_id"`
	Type        string    `db:"type" json:"type"`
	Phone       *string   `db:"phone" json:"phone,omitempty"`
	Line1       string    `db:"line1" json:"line1"`
	Line2       *string   `db:"line2" json:"line2,omitempty"`
	City        string    `db:"city" json:"city"`
	Region      *string   `db:"region" json:"region,omitempty"`
	PostalCode  string    `db:"postal_code" json:"postal_code"`
	CountryCode string    `db:"country_code" json:"country_code"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
