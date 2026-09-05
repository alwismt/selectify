package model

import (
	"time"
)

type Products []Product

type Product struct {
	ID          uint         `db:"product_id" json:"productId"`
	MerchantID  *uint64      `db:"merchant_id" json:"merchantId"`
	SKU         string       `db:"sku" json:"sku"`
	Name        string       `db:"name" json:"name"`
	Description *string      `db:"description" json:"description"`
	Slug        string       `db:"slug" json:"slug"`
	PriceAmount uint         `db:"price_amount" json:"price"`
	IsActive    bool         `db:"is_active" json:"isActive"`
	InStock     bool         `db:"in_stock" json:"inStock"`
	CreatedAt   time.Time    `db:"created_at" json:"-"`
	UpdatedAt   time.Time    `db:"updated_at" json:"-"`
	DeletedAt   *time.Time   `db:"deleted_at" json:"-"`
	ProductFile *ProductFile `json:"productFile,omitempty"`
	Categories  Categories   `json:"category,omitempty"`
}
