package catalog

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type fakeRepo struct {
	listFn func(ctx context.Context, limit, offset int, search string) ([]ProductListItem, error)
	getFn  func(ctx context.Context, slug string) (*ProductDetail, error)
}

func (f fakeRepo) ListProducts(ctx context.Context, limit, offset int, search string) ([]ProductListItem, error) {
	return f.listFn(ctx, limit, offset, search)
}
func (f fakeRepo) GetProductBySlug(ctx context.Context, slug string) (*ProductDetail, error) {
	return f.getFn(ctx, slug)
}

func TestCatalog_ListProducts_OK(t *testing.T) {
	repo := fakeRepo{
		listFn: func(ctx context.Context, limit, offset int, search string) ([]ProductListItem, error) {
			return []ProductListItem{
				{ID: "1", Slug: "abc", Name: "ABC", MinPrice: 1000, MaxPrice: 2000, ImageURL: "x"},
			}, nil
		},
		getFn: func(ctx context.Context, slug string) (*ProductDetail, error) {
			return nil, ErrNotFound
		},
	}

	svc := NewService(repo)
	h := NewHandler(svc)

	r := chi.NewRouter()
	h.Routes(r)

	req := httptest.NewRequest(http.MethodGet, "/products?limit=10&offset=0", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.NotNil(t, body["items"])
}

func TestCatalog_GetProduct_404(t *testing.T) {
	repo := fakeRepo{
		listFn: func(ctx context.Context, limit, offset int, search string) ([]ProductListItem, error) {
			return nil, nil
		},
		getFn: func(ctx context.Context, slug string) (*ProductDetail, error) {
			return nil, ErrNotFound
		},
	}

	svc := NewService(repo)
	h := NewHandler(svc)

	r := chi.NewRouter()
	h.Routes(r)

	req := httptest.NewRequest(http.MethodGet, "/products/not-exist", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}
