package model

import (
	"time"

	"github.com/google/uuid"
)

type ProductFiles []ProductFile

type ProductFile struct {
	FileID      uuid.UUID `db:"product_file_id" json:"file_id"`
	ProductID   uint      `db:"product_id" json:"product_id"`
	VariantID   *uint     `db:"variant_id" json:"variant_id,omitempty"`
	ContentType string    `db:"content_type" json:"content_type"`
	Position    uint      `db:"position" json:"position"`
	IsPrimary   bool      `db:"is_primary" json:"is_primary"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
}
