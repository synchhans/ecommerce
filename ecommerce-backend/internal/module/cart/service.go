package cart

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCart(ctx context.Context) (string, error) {
	return s.repo.CreateCart(ctx)
}

func (s *Service) GetCart(ctx context.Context, cartID string) (*Cart, error) {
	return s.repo.GetCart(ctx, cartID)
}

func (s *Service) AddOrReplaceItem(ctx context.Context, cartID, variantID string, qty int) error {
	return s.repo.UpsertItem(ctx, cartID, variantID, qty)
}

func (s *Service) UpdateItemQty(ctx context.Context, cartID, itemID string, qty int) error {
	return s.repo.UpdateItemQty(ctx, cartID, itemID, qty)
}

func (s *Service) RemoveItem(ctx context.Context, cartID, itemID string) error {
	return s.repo.DeleteItem(ctx, cartID, itemID)
}
