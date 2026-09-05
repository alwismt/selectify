package request

type ProductVariantRequest struct {
	SKU        string                    `json:"sku" validate:"required"`
	Price      uint64                    `json:"price" validate:"required,gt=0"`
	Quantity   uint                      `json:"quantity" validate:"required,gt=0"`
	Attributes []VariantAttributeRequest `json:"attributes" validate:"dive"`
}

type VariantAttributeRequest struct {
	Name  string `json:"name" validate:"required"`
	Value string `json:"value" validate:"required"`
}
