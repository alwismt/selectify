package types

type UserRole string

const (
	RoleCustomer      UserRole = "customer"
	RoleMerchant      UserRole = "merchant"
	RolePlatformAdmin UserRole = "platform_admin"
	RoleVendorAdmin   UserRole = "vendor_admin"
	RoleSuperAdmin    UserRole = "super_admin"
)

func (r UserRole) Valid() bool {
	switch r {
	case RoleCustomer,
		RoleMerchant,
		RolePlatformAdmin,
		RoleVendorAdmin,
		RoleSuperAdmin:
		return true
	default:
		return false
	}
}
