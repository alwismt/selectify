package types

type MerchantRole string

const (
	MerchantRoleOwner   MerchantRole = "owner"
	MerchantRoleAdmin   MerchantRole = "admin"
	MerchantRoleManager MerchantRole = "manager"
	MerchantRoleStaff   MerchantRole = "staff"
)

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
