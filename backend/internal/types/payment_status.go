package types

type PaymentStatus string

var (
	PaymentStatusPending = PaymentStatus("pending")
	PaymentStatusSuccess = PaymentStatus("success")
)
