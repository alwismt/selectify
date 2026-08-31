package types

type MerchantStatus string

const (
	MerchantStatusActive   MerchantStatus = "active"
	MerchantStatusPending  MerchantStatus = "pending"
	MerchantStatusRejected MerchantStatus = "rejected"
	MerchantStatusClosed   MerchantStatus = "closed"
)

type MerchantVerificationStatus string

const (
	MerchantVerificationStatusPending  MerchantVerificationStatus = "pending"
	MerchantVerificationStatusRejected MerchantVerificationStatus = "rejected"
	MerchantVerificationStatusVerified MerchantVerificationStatus = "verified"
)
