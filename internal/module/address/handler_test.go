package address

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	httpx "github.com/synchhans/ecommerce-backend/internal/platform/http"
)

type fakeRepo struct {
	createFn func(ctx context.Context, userID string, a Address) (*Address, error)
	listFn   func(ctx context.Context, userID string) ([]Address, error)
	updateFn func(ctx context.Context, userID, addressID string, a Address) (*Address, error)
	deleteFn func(ctx context.Context, userID, addressID string) error
	defFn    func(ctx context.Context, userID, addressID string) error
}

func (f fakeRepo) Create(ctx context.Context, userID string, a Address) (*Address, error) {
	return f.createFn(ctx, userID, a)
}
func (f fakeRepo) List(ctx context.Context, userID string) ([]Address, error) {
	return f.listFn(ctx, userID)
}
func (f fakeRepo) Update(ctx context.Context, userID, addressID string, a Address) (*Address, error) {
	return f.updateFn(ctx, userID, addressID, a)
}
func (f fakeRepo) Delete(ctx context.Context, userID, addressID string) error {
	return f.deleteFn(ctx, userID, addressID)
}
func (f fakeRepo) SetDefault(ctx context.Context, userID, addressID string) error {
	return f.defFn(ctx, userID, addressID)
}

func TestAddresses_Unauthorized(t *testing.T) {
	repo := fakeRepo{
		createFn: func(ctx context.Context, userID string, a Address) (*Address, error) { return nil, nil },
		listFn:   func(ctx context.Context, userID string) ([]Address, error) { return nil, nil },
		updateFn: func(ctx context.Context, userID, addressID string, a Address) (*Address, error) { return nil, nil },
		deleteFn: func(ctx context.Context, userID, addressID string) error { return nil },
		defFn:    func(ctx context.Context, userID, addressID string) error { return nil },
	}
	svc := NewService(repo)
	h := NewHandler(svc)

	r := chi.NewRouter()
	h.Routes(r)

	req := httptest.NewRequest(http.MethodGet, "/addresses", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAddresses_Create_201(t *testing.T) {
	repo := fakeRepo{
		createFn: func(ctx context.Context, userID string, a Address) (*Address, error) {
			require.Equal(t, "u1", userID)
			return &Address{ID: "a1", RecipientName: a.RecipientName, City: a.City, Country: "ID"}, nil
		},
		listFn:   func(ctx context.Context, userID string) ([]Address, error) { return nil, nil },
		updateFn: func(ctx context.Context, userID, addressID string, a Address) (*Address, error) { return nil, nil },
		deleteFn: func(ctx context.Context, userID, addressID string) error { return nil },
		defFn:    func(ctx context.Context, userID, addressID string) error { return nil },
	}
	svc := NewService(repo)
	h := NewHandler(svc)

	secret := []byte("secret")
	token, err := httpx.SignJWT("u1", secret, 1*time.Hour)
	require.NoError(t, err)

	r := chi.NewRouter()
	r.Use(httpx.AuthMiddleware(secret))
	h.Routes(r)

	body := []byte(`{"recipient_name":"A","phone":"1","address_line1":"x","city":"y","province":"z","postal_code":"1","country":"ID"}`)
	req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var out Address
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &out))
	require.Equal(t, "a1", out.ID)
}
