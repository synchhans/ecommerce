package catalog

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListProducts(ctx context.Context, limit, offset int, search string) ([]ProductListItem, error) {
	return s.repo.ListProducts(ctx, limit, offset, search)
}

func (s *Service) GetProductDetail(ctx context.Context, slug string) (*ProductDetail, error) {
	if slug == "" {
		return nil, ErrNotFound
	}
	p, err := s.repo.GetProductBySlug(ctx, slug)
	if err != nil {
		// repository bisa return pgx.ErrNoRows; handler akan map ke 404 via ErrNotFound
		return nil, err
	}
	return p, nil
}
