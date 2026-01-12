package order

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
	createFn func(ctx context.Context, cartID string, addr AddressSnapshot) (string, error)
	getFn    func(ctx context.Context, orderID string) (*Order, error)
}

func (f fakeRepo) CreateOrderFromCart(ctx context.Context, cartID string, addr AddressSnapshot) (string, error) {
	return f.createFn(ctx, cartID, addr)
}
func (f fakeRepo) GetOrder(ctx context.Context, orderID string) (*Order, error) { return f.getFn(ctx, orderID) }

func TestCheckout_201(t *testing.T) {
	repo := fakeRepo{
		createFn: func(ctx context.Context, cartID string, addr AddressSnapshot) (string, error) {
			require.Equal(t, "cart-1", cartID)
			return "order-1", nil
		},
		getFn: func(ctx context.Context, orderID string) (*Order, error) { return nil, nil },
	}
	svc := NewService(repo)
	h := NewHandler(svc)

	r := chi.NewRouter()
	h.Routes(r)

	reqBody := []byte(`{"cart_id":"cart-1","address":{"recipient_name":"A","phone":"1","address_line1":"x","city":"y","province":"z","postal_code":"1","country":"ID"}}`)
	req := httptest.NewRequest(http.MethodPost, "/checkout", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, "order-1", body["order_id"])
}

func TestCheckout_400_InvalidJSON(t *testing.T) {
	repo := fakeRepo{
		createFn: func(ctx context.Context, cartID string, addr AddressSnapshot) (string, error) { return "", nil },
		getFn:    func(ctx context.Context, orderID string) (*Order, error) { return nil, nil },
	}
	svc := NewService(repo)
	h := NewHandler(svc)
	r := chi.NewRouter()
	h.Routes(r)

	req := httptest.NewRequest(http.MethodPost, "/checkout", bytes.NewReader([]byte(`{`)))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetOrder_404(t *testing.T) {
	repo := fakeRepo{
		createFn: func(ctx context.Context, cartID string, addr AddressSnapshot) (string, error) { return "", nil },
		getFn: func(ctx context.Context, orderID string) (*Order, error) {
			return nil, ErrNotFound
		},
	}
	svc := NewService(repo)
	h := NewHandler(svc)
	r := chi.NewRouter()
	h.Routes(r)

	req := httptest.NewRequest(http.MethodGet, "/orders/x", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}
