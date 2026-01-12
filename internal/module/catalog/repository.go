package catalog

import "context"

type Repository interface {
	ListProducts(ctx context.Context, limit, offset int, search string) ([]ProductListItem, error)
	GetProductBySlug(ctx context.Context, slug string) (*ProductDetail, error)
}
