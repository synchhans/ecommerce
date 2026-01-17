package cart

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type fakeRepo struct {
	createFn func(ctx context.Context) (string, error)
	getFn    func(ctx context.Context, cartID string) (*Cart, error)
	upsertFn func(ctx context.Context, cartID, variantID string, qty int) error
	updateFn func(ctx context.Context, cartID, itemID string, qty int) error
	deleteFn func(ctx context.Context, cartID, itemID string) error
}

func (f fakeRepo) CreateCart(ctx context.Context) (string, error) { return f.createFn(ctx) }
func (f fakeRepo) GetCart(ctx context.Context, cartID string) (*Cart, error) {
	return f.getFn(ctx, cartID)
}
func (f fakeRepo) UpsertItem(ctx context.Context, cartID, variantID string, qty int) error {
	return f.upsertFn(ctx, cartID, variantID, qty)
}
func (f fakeRepo) UpdateItemQty(ctx context.Context, cartID, itemID string, qty int) error {
	return f.updateFn(ctx, cartID, itemID, qty)
}
func (f fakeRepo) DeleteItem(ctx context.Context, cartID, itemID string) error {
	return f.deleteFn(ctx, cartID, itemID)
}

func TestCart_CreateCart_201(t *testing.T) {
	repo := fakeRepo{
		createFn: func(ctx context.Context) (string, error) { return "cart-1", nil },
		getFn:    func(ctx context.Context, cartID string) (*Cart, error) { return nil, nil },
		upsertFn: func(ctx context.Context, cartID, variantID string, qty int) error { return nil },
		updateFn: func(ctx context.Context, cartID, itemID string, qty int) error { return nil },
		deleteFn: func(ctx context.Context, cartID, itemID string) error { return nil },
	}
	svc := NewService(repo)
	h := NewHandler(svc)

	r := chi.NewRouter()
	h.Routes(r)

	req := httptest.NewRequest(http.MethodPost, "/cart", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, "cart-1", body["id"])
}

func TestCart_GetCart_404(t *testing.T) {
	repo := fakeRepo{
		createFn: func(ctx context.Context) (string, error) { return "", nil },
		getFn: func(ctx context.Context, cartID string) (*Cart, error) {
			return nil, ErrNotFound
		},
		upsertFn: func(ctx context.Context, cartID, variantID string, qty int) error { return nil },
		updateFn: func(ctx context.Context, cartID, itemID string, qty int) error { return nil },
		deleteFn: func(ctx context.Context, cartID, itemID string) error { return nil },
	}
	svc := NewService(repo)
	h := NewHandler(svc)

	r := chi.NewRouter()
	h.Routes(r)

	req := httptest.NewRequest(http.MethodGet, "/cart/x", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCart_UpsertItem_400_InvalidPayload(t *testing.T) {
	repo := fakeRepo{
		createFn: func(ctx context.Context) (string, error) { return "", nil },
		getFn:    func(ctx context.Context, cartID string) (*Cart, error) { return nil, nil },
		upsertFn: func(ctx context.Context, cartID, variantID string, qty int) error { return nil },
		updateFn: func(ctx context.Context, cartID, itemID string, qty int) error { return nil },
		deleteFn: func(ctx context.Context, cartID, itemID string) error { return nil },
	}
	svc := NewService(repo)
	h := NewHandler(svc)

	r := chi.NewRouter()
	h.Routes(r)

	reqBody := []byte(`{"variant_id":"","qty":0}`)
	req := httptest.NewRequest(http.MethodPost, "/cart/1/items", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
