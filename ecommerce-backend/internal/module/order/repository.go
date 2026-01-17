package order

import "context"

type Repository interface {
	CreateOrderFromCart(ctx context.Context, cartID string, shipAddr AddressSnapshot) (string, error)
	GetOrder(ctx context.Context, orderID string) (*Order, error)
}
