package request

type OrderShippingAddressReq struct {
	Phone       *string `json:"phone" validate:"omitempty,max=50"`
	Line1       string  `json:"line1" validate:"required,min=1,max=255"`
	Line2       *string `json:"line2" validate:"omitempty,max=255"`
	City        string  `json:"city" validate:"required,min=1,max=100"`
	Region      *string `json:"region" validate:"omitempty,max=100"`
	PostalCode  string  `json:"postal_code" validate:"required,min=1,max=32"`
	CountryCode string  `json:"country_code" validate:"required,len=2"`
}
