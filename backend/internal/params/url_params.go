package params

import "fmt"

const (
	ProductId   = "product_id"
	ProductSlug = "product_slug"
	CartItemId  = "item_id"
	OrderId     = "order_id"
)

var (
	ProductIdParm   = fmt.Sprintf("{%s}", ProductId)
	ProductSlugParm = fmt.Sprintf("{%s}", ProductSlug)
	CartItemIdParam = fmt.Sprintf("{%s}", CartItemId)
	OrderIdParam    = fmt.Sprintf("{%s}", OrderId)
)
