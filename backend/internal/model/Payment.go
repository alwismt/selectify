package model

import "time"

type Payment struct {
	ID      uint64 `db:"id" json:"id"`
	OrderID uint64 `db:"order_id" json:"orderId"`

	Provider          string `db:"provider" json:"provider"`
	ProviderPaymentID string `db:"provider_payment_id" json:"providerPaymentId"`
	ClientSecret      string `db:"client_secret" json:"clientSecret,omitempty"`

	Status string `db:"status" json:"status"`

	Amount   int64  `db:"amount" json:"amount"`
	Currency string `db:"currency" json:"currency"`

	FailureCode    *string `db:"failure_code" json:"failureCode,omitempty"`
	FailureMessage *string `db:"failure_message" json:"failureMessage,omitempty"`

	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
	PaidAt    *time.Time `db:"paid_at" json:"paidAt,omitempty"`
}
