package types

type UserStatus string

var (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "disabled"
	UserStatusPending  UserStatus = "pending"
)
