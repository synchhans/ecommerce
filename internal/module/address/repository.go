package address

import "context"

type Repository interface {
	Create(ctx context.Context, userID string, a Address) (*Address, error)
	List(ctx context.Context, userID string) ([]Address, error)
	Update(ctx context.Context, userID, addressID string, a Address) (*Address, error)
	Delete(ctx context.Context, userID, addressID string) error
	SetDefault(ctx context.Context, userID, addressID string) error
}
