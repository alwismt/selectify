package model

import "time"

type Currency struct {
	Code      string    `db:"code" json:"code"`
	Name      string    `db:"name" json:"name"`
	MinorUnit uint      `db:"minor_unit" json:"minorUnit"`
	IsActive  bool      `db:"is_active" json:"isActive"`
	IsDefault bool      `db:"is_default" json:"isDefault"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}
