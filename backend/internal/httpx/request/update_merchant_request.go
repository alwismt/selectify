package request

import "strings"

type UpdateMerchantRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty"`
	CountryCode *string `json:"countryCode,omitempty" validate:"omitempty,len=2"`
}

func (r *UpdateMerchantRequest) Empty() bool {
	return r.Name == nil &&
		r.Description == nil &&
		r.CountryCode == nil
}

func (r *UpdateMerchantRequest) Sanitize() {
	if r.Name != nil {
		v := strings.TrimSpace(*r.Name)
		r.Name = &v
	}
	if r.Description != nil {
		v := strings.TrimSpace(*r.Description)
		r.Description = &v
	}
	if r.CountryCode != nil {
		v := strings.ToUpper(strings.TrimSpace(*r.CountryCode))
		r.CountryCode = &v
	}
}
