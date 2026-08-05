package types

type UserRole string

const (
	RoleCustomer      UserRole = "customer"
	RolePlatformAdmin UserRole = "platform_admin"
	RoleVendorAdmin   UserRole = "vendor_admin"
	RoleSuperAdmin    UserRole = "super_admin"
)
