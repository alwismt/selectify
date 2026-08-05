package types

type OrderStatus string

// -- pending|paid|cancelled|shipped|refunded

var (
	OrderStatusNew       OrderStatus = "new"
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusFailed    OrderStatus = "failed"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusRefunded  OrderStatus = "refunded"
	OrderStatusShipped   OrderStatus = "shipped"
)
