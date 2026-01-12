package catalog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService_GetProductDetail_EmptySlug_NotFound(t *testing.T) {
	repo := fakeRepo{
		listFn: func(ctx context.Context, limit, offset int, search string) ([]ProductListItem, error) {
			return nil, nil
		},
		getFn: func(ctx context.Context, slug string) (*ProductDetail, error) {
			return nil, ErrNotFound
		},
	}

	svc := NewService(repo)
	p, err := svc.GetProductDetail(context.Background(), "")
	require.Error(t, err)
	require.Nil(t, p)
}
