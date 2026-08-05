package response

type CartResponse struct {
	Items     []CartItem `json:"items"`
	Currency  string     `json:"currency"`
	Subtotal  float64    `json:"subtotal"`
	ItemCount uint       `json:"item_count"`
}

type CartItem struct {
	ID       uint           `json:"id"`
	Quantity uint           `json:"quantity"`
	Variant  ProductVariant `json:"variant"`
}

type ProductVariant struct {
	ID          uint              `json:"id"`
	SKU         string            `json:"sku"`
	PriceAmount *float64          `json:"price_amount,omitempty"`
	Currency    string            `json:"currency"`
	Attributes  map[string]string `json:"attributes"`
	Product     Product           `json:"product"`
	StockQty    uint              `json:"available_qty"`
}

type Product struct {
	ID          uint    `json:"id"`
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description"`
}
