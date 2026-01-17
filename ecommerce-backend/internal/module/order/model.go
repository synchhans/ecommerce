package order

type Order struct {
	ID          string      `json:"id"`
	OrderNumber string      `json:"order_number"`
	Status      string      `json:"status"`
	Currency    string      `json:"currency"`
	Subtotal    int64       `json:"subtotal"`
	Discount    int64       `json:"discount_total"`
	Shipping    int64       `json:"shipping_total"`
	GrandTotal  int64       `json:"grand_total"`
	Items       []OrderItem `json:"items"`
}

type OrderItem struct {
	ID        string `json:"id"`
	VariantID string `json:"variant_id"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	UnitPrice int64  `json:"unit_price"`
	Qty       int    `json:"qty"`
	LineTotal int64  `json:"line_total"`
}

type AddressSnapshot struct {
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	AddressLine1  string `json:"address_line1"`
	AddressLine2  string `json:"address_line2,omitempty"`
	City          string `json:"city"`
	Province      string `json:"province"`
	PostalCode    string `json:"postal_code"`
	Country       string `json:"country"`
}
