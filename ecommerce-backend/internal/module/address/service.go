package address

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) Create(ctx context.Context, userID string, a Address) (*Address, error) {
	return s.repo.Create(ctx, userID, a)
}
func (s *Service) List(ctx context.Context, userID string) ([]Address, error) {
	return s.repo.List(ctx, userID)
}
func (s *Service) Update(ctx context.Context, userID, addressID string, a Address) (*Address, error) {
	return s.repo.Update(ctx, userID, addressID, a)
}
func (s *Service) Delete(ctx context.Context, userID, addressID string) error {
	return s.repo.Delete(ctx, userID, addressID)
}
func (s *Service) SetDefault(ctx context.Context, userID, addressID string) error {
	return s.repo.SetDefault(ctx, userID, addressID)
}
