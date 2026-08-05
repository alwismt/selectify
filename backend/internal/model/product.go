package model

import (
	"time"
)

type Products []Product

type Product struct {
	ID          uint         `db:"product_id" json:"productId"`
	SKU         string       `db:"sku" json:"sku"`
	Name        string       `db:"name" json:"name"`
	Description *string      `db:"description" json:"description"`
	Slug        *string      `db:"slug" json:"slug,omitempty"`
	PriceAmount float64      `db:"price_amount" json:"priceAmount"`
	Currency    string       `db:"currency" json:"currency"`
	IsActive    bool         `db:"is_active" json:"isActive"`
	InStock     bool         `db:"in_stock" json:"inStock"`
	CreatedAt   time.Time    `db:"created_at" json:"-"`
	UpdatedAt   time.Time    `db:"updated_at" json:"-"`
	DeletedAt   *time.Time   `db:"deleted_at" json:"-"`
	ProductFile *ProductFile `json:"productFile,omitempty"`
}
