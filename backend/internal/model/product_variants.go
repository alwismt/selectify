package model

import (
	"context"
	"encoding/json"
	"time"

	"alwis.dev/selectify/internal/logger"
)

type ProductVariants []ProductVariant

type ProductVariant struct {
	ID                       uint                     `db:"id" json:"id"`
	ProductID                uint                     `db:"product_id" json:"product_id"`
	SKU                      string                   `db:"sku" json:"sku"`
	PriceAmount              *float64                 `db:"price_amount" json:"price_amount,omitempty"`
	Currency                 string                   `db:"currency" json:"currency"`
	IsActive                 bool                     `db:"is_active" json:"is_active"`
	CreatedAt                time.Time                `db:"created_at" json:"created_at"`
	UpdatedAt                time.Time                `db:"updated_at" json:"updated_at"`
	DeletedAt                *time.Time               `db:"deleted_at" json:"deleted_at,omitempty"`
	StockQty                 uint                     `db:"stock_qty" json:"stock_quantity"`
	ReservedQty              uint                     `db:"reserved_qty" json:"reserved_quantity"`
	ProductVariantAttributes ProductVariantAttributes `db:"-" json:"product_variant_attributes"`
	ProductFiles             ProductFiles             `db:"-" json:"files"`
}

type ProductVariantAttributes []ProductVariantAttribute

type ProductVariantAttribute struct {
	ID        uint   `db:"id" json:"id"`
	VariantID uint   `db:"variant_id" json:"variant_id"`
	Name      string `db:"name" json:"name"`
	Value     string `db:"value" json:"value"`
}

func (i *ProductVariant) AvailableStockQty() uint {
	return i.StockQty - i.ReservedQty
}

func (pva ProductVariantAttributes) ToJson(ctx context.Context) ([]byte, error) {
	attrsJSON, err := json.Marshal(pva)
	if err != nil {
		return nil, logger.Errorf(ctx, err, "failed to marshal variant attributes")
	}

	return attrsJSON, nil
}
