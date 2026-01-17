package catalog

type ProductListItem struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	MinPrice    int64  `json:"min_price"`
	MaxPrice    int64  `json:"max_price"`
	ImageURL    string `json:"image_url,omitempty"`
	IsActive    bool   `json:"-"`
}

type ProductDetail struct {
	ID          string          `json:"id"`
	Slug        string          `json:"slug"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Images      []ProductImage  `json:"images"`
	Variants    []ProductVariant `json:"variants"`
}

type ProductImage struct {
	URL      string `json:"url"`
	Position int    `json:"position"`
}

type ProductVariant struct {
	ID     string `json:"id"`
	SKU    string `json:"sku"`
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	Active bool   `json:"active"`
}
