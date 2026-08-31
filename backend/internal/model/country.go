package model

type Countries []Country

type Country struct {
	Code      string `db:"code" json:"code"`
	Name      string `db:"name" json:"name"`
	IsActive  bool   `db:"is_active" json:"-"`
	CreatedAt string `db:"created_at" json:"-"`
	UpdatedAt string `db:"updated_at" json:"-"`
}
