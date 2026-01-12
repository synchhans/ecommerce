package inventory

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
	getFn func(ctx context.Context, variantID string) (*Availability, error)
}

func (f fakeRepo) GetAvailability(ctx context.Context, variantID string) (*Availability, error) {
	return f.getFn(ctx, variantID)
}

func TestAvailability_200(t *testing.T) {
	repo := fakeRepo{
		getFn: func(ctx context.Context, variantID string) (*Availability, error) {
			require.Equal(t, "v-1", variantID)
			return &Availability{
				VariantID:   "v-1",
				StockOnHand: 10,
				Reserved:    3,
				Available:   7,
			}, nil
		},
	}
	svc := NewService(repo)
	h := NewHandler(svc)

	r := chi.NewRouter()
	h.Routes(r)

	req := httptest.NewRequest(http.MethodGet, "/variants/v-1/availability", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var out Availability
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &out))
	require.Equal(t, 7, out.Available)
}

func TestAvailability_404(t *testing.T) {
	repo := fakeRepo{
		getFn: func(ctx context.Context, variantID string) (*Availability, error) {
			return nil, ErrNotFound
		},
	}
	svc := NewService(repo)
	h := NewHandler(svc)

	r := chi.NewRouter()
	h.Routes(r)

	req := httptest.NewRequest(http.MethodGet, "/variants/x/availability", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}
