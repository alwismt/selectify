package permission

import "alwis.dev/selectify/internal/types"

// Merchant permissions
const (
	MerchantRead   Permission = "merchant_read"
	MerchantUpdate Permission = "merchant_update"
)

// Product permissions
const (
	ProductCreate  Permission = "product_create"
	ProductRead    Permission = "product_read"
	ProductUpdate  Permission = "product_update"
	ProductDelete  Permission = "product_delete"
	ProductPublish Permission = "product_publish"
)

// Grouped permission objects
type permissionObject struct {
	Create Permission
	Read   Permission
	Update Permission
	Delete Permission
}

var Merchant = struct {
	Read   Permission
	Update Permission
}{
	Read:   MerchantRead,
	Update: MerchantUpdate,
}

var Product = struct {
	permissionObject
	Publish Permission
}{
	permissionObject: permissionObject{
		Create: ProductCreate,
		Read:   ProductRead,
		Update: ProductUpdate,
		Delete: ProductDelete,
	},
	Publish: ProductPublish,
}

// Map merchant roles directly to permissions
var merchantRolePermissions = map[types.MerchantRolePermission][]Permission{
	types.MerchantRoleOwnerPermission: {
		Merchant.Read,
		Merchant.Update,

		Product.Create,
		Product.Read,
		Product.Update,
		Product.Delete,
		Product.Publish,
	},

	types.MerchantRoleAdminPermission: {
		Merchant.Read,
		Merchant.Update,

		Product.Create,
		Product.Read,
		Product.Update,
		Product.Delete,
		Product.Publish,
	},

	types.MerchantRoleManagerPermission: {
		Merchant.Read,

		Product.Create,
		Product.Read,
		Product.Update,
		Product.Publish,
	},

	types.MerchantRoleStaffPermission: {
		Merchant.Read,
		Product.Read,
	},
}

// HasMerchantPermission checks if a merchant role has the required permissions
func HasMerchantPermission(role types.MerchantRolePermission, required ...Permission) bool {
	actual := merchantRolePermissions[role]
	return hasPermission(actual, required...)
}

// hasPermission checks if any of the required permissions exists in a permission slice
func hasPermission(actualPermissions []Permission, requiredPermissions ...Permission) bool {
	if len(requiredPermissions) == 0 {
		return true
	}

	for _, requiredPermission := range requiredPermissions {
		for _, actualPermission := range actualPermissions {
			if actualPermission == requiredPermission {
				return true
			}
		}
	}

	return false
}

// GetMerchantPermissions returns a copy of all permissions for a merchant role
func GetMerchantPermissions(role types.MerchantRolePermission) []Permission {
	perms, ok := merchantRolePermissions[role]
	if !ok {
		return nil
	}

	return append([]Permission(nil), perms...)
}
