package model

import "time"

type Categories []*Category

type Category struct {
	CategoryID uint       `db:"category_id" json:"categoryId"`
	Name       string     `db:"name" json:"name"`
	Slug       string     `db:"slug" json:"slug"`
	ParentID   *uint      `db:"parent_id" json:"parentId,omitempty"`
	IsActive   bool       `db:"is_active" json:"isActive"`
	CreatedAt  time.Time  `db:"created_at" json:"-"`
	UpdatedAt  time.Time  `db:"updated_at" json:"-"`
	Children   Categories `json:"children,omitempty"`
}
