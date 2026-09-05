package model

import (
	"time"

	"alwis.dev/selectify/internal/types"
)

type Merchant struct {
	MerchantID  uint64  `db:"merchant_id" json:"merchantId"`
	Name        string  `db:"name" json:"name"`
	Slug        string  `db:"slug" json:"slug"`
	Description *string `db:"description" json:"description,omitempty"`
	Logo        *string `db:"logo" json:"logo,omitempty"`
	CountryCode string  `db:"country_code" json:"countryCode"`

	Status                   types.MerchantStatus             `db:"status" json:"status"`
	VerificationStatus       types.MerchantVerificationStatus `db:"verification_status" json:"verificationStatus"`
	PaymentProvider          *types.PaymentProvider           `db:"payment_provider" json:"paymentProvider,omitempty"`
	PaymentProviderAccountID *string                          `db:"payment_provider_account_id" json:"paymentProviderAccountId,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
