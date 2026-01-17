package payment

import "context"

type Repository interface {
	InitiatePayment(ctx context.Context, orderID, provider string) (*InitiateResult, error)
	HandleWebhook(ctx context.Context, provider, providerRef, newStatus string, rawPayload []byte) (*WebhookResult, error)
}
