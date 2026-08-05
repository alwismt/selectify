package request

type AddCartReq struct {
	VariantId uint `json:"variant_id" validate:"required,gt=0"`
	Quantity  uint `json:"quantity" validate:"required,gte=1"`
}

type QntReq struct {
	Quantity uint `json:"quantity" validate:"required,gte=1"`
}
