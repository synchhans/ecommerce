package payment

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
	initFn func(ctx context.Context, orderID, provider string) (*InitiateResult, error)
	whFn   func(ctx context.Context, provider, providerRef, status string, raw []byte) (*WebhookResult, error)
}

func (f fakeRepo) InitiatePayment(ctx context.Context, orderID, provider string) (*InitiateResult, error) {
	return f.initFn(ctx, orderID, provider)
}
func (f fakeRepo) HandleWebhook(ctx context.Context, provider, providerRef, newStatus string, rawPayload []byte) (*WebhookResult, error) {
	return f.whFn(ctx, provider, providerRef, newStatus, rawPayload)
}

func TestInitiate_201(t *testing.T) {
	repo := fakeRepo{
		initFn: func(ctx context.Context, orderID, provider string) (*InitiateResult, error) {
			require.Equal(t, "order-1", orderID)
			require.Equal(t, "manual", provider)
			return &InitiateResult{
				PaymentID:   "pay-1",
				OrderID:     "order-1",
				Status:      "initiated",
				Amount:      1000,
				Provider:    "manual",
				ProviderRef: "ref-1",
			}, nil
		},
		whFn: func(ctx context.Context, provider, providerRef, status string, raw []byte) (*WebhookResult, error) {
			return nil, nil
		},
	}

	svc := NewService(repo)
	h := NewHandler(svc)
	r := chi.NewRouter()
	h.Routes(r)

	reqBody := []byte(`{"order_id":"order-1","provider":"manual"}`)
	req := httptest.NewRequest(http.MethodPost, "/payments/initiate", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var out InitiateResult
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &out))
	require.Equal(t, "pay-1", out.PaymentID)
	require.Equal(t, "order-1", out.OrderID)
}

func TestWebhook_400_InvalidStatus(t *testing.T) {
	repo := fakeRepo{
		initFn: func(ctx context.Context, orderID, provider string) (*InitiateResult, error) { return nil, nil },
		whFn: func(ctx context.Context, provider, providerRef, status string, raw []byte) (*WebhookResult, error) {
			return nil, ErrInvalidStatus
		},
	}

	svc := NewService(repo)
	h := NewHandler(svc)
	r := chi.NewRouter()
	h.Routes(r)

	reqBody := []byte(`{"provider_ref":"ref-1","status":"weird"}`)
	req := httptest.NewRequest(http.MethodPost, "/payments/webhook/manual", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
