package inventory

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Availability(ctx context.Context, variantID string) (*Availability, error) {
	return s.repo.GetAvailability(ctx, variantID)
}
