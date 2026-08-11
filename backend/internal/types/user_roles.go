package types

type UserRole string

const (
	RoleCustomer   UserRole = "customer"
	RoleMerchant   UserRole = "merchant"
	RoleSuperAdmin UserRole = "super_admin"
)

func (r UserRole) Valid() bool {
	switch r {
	case RoleCustomer,
		RoleMerchant,
		RoleSuperAdmin:
		return true
	default:
		return false
	}
}

func (r UserRole) IsMerchant() bool {
	return r == RoleMerchant
}
