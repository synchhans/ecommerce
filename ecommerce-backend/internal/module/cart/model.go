package cart

type Cart struct {
	ID     string     `json:"id"`
	Status string     `json:"status"`
	Items  []CartItem `json:"items"`
}

type CartItem struct {
	ID        string `json:"id"`
	VariantID string `json:"variant_id"`
	Qty       int    `json:"qty"`
}
