package request

import (
	"strings"

	"alwis.dev/selectify/internal/model"
)

type ProductRequest struct {
	Name        string           `json:"name" validate:"required"`
	Description string           `json:"description" validate:"required"`
	Price       uint             `json:"price" validate:"required,gt=0"`
	SKU         string           `json:"sku" validate:"required,sku"`
	CategoryIDs []uint           `json:"categoryIds" validate:"required,min=1,dive,gt=0"`
	Image       model.FileStream `json:"image"`
}

func (r *ProductRequest) Sanitize() {
	if r.Name != "" {
		r.Name = strings.TrimSpace(r.Name)
	}
	if r.Description != "" {
		r.Description = strings.TrimSpace(r.Description)
	}
	if r.SKU != "" {
		r.SKU = strings.ToUpper(strings.TrimSpace(r.SKU))
	}
}
