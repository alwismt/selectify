package types

import (
	"slices"
)

type MerchantRolePermission string

const (
	MerchantRoleOwnerPermission   MerchantRolePermission = "owner"
	MerchantRoleAdminPermission   MerchantRolePermission = "admin"
	MerchantRoleManagerPermission MerchantRolePermission = "manager"
	MerchantRoleStaffPermission   MerchantRolePermission = "staff"
)

type MerchantRolePermissions []MerchantRolePermission

func (ps MerchantRolePermissions) HasPermission(role MerchantRolePermission) bool {
	return slices.Contains(ps, role)
}
