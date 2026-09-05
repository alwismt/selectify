package params

import "fmt"

const (
	ProductId   = "product_id"
	ProductPath = "product_path"
	CartItemId  = "item_id"
	OrderId     = "order_id"
)

var (
	ProductIdParm   = fmt.Sprintf("{%s}", ProductId)
	ProductPathParm = fmt.Sprintf("{%s}", ProductPath)
	CartItemIdParam = fmt.Sprintf("{%s}", CartItemId)
	OrderIdParam    = fmt.Sprintf("{%s}", OrderId)
)
