package cart

import "context"

type Repository interface {
	CreateCart(ctx context.Context) (string, error)
	GetCart(ctx context.Context, cartID string) (*Cart, error)

	UpsertItem(ctx context.Context, cartID, variantID string, qty int) error
	UpdateItemQty(ctx context.Context, cartID, itemID string, qty int) error
	DeleteItem(ctx context.Context, cartID, itemID string) error
}
