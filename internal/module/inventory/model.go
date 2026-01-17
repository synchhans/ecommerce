package inventory

type Availability struct {
	VariantID    string `json:"variant_id"`
	StockOnHand  int    `json:"stock_on_hand"`
	Reserved     int    `json:"reserved"`
	Available    int    `json:"available"`
}
