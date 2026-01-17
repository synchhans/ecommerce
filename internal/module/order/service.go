package order

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Checkout(ctx context.Context, cartID string, addr AddressSnapshot) (string, error) {
	return s.repo.CreateOrderFromCart(ctx, cartID, addr)
}

func (s *Service) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	return s.repo.GetOrder(ctx, orderID)
}
