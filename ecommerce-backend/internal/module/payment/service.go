package payment

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Initiate(ctx context.Context, orderID, provider string) (*InitiateResult, error) {
	if provider == "" {
		provider = "manual"
	}
	return s.repo.InitiatePayment(ctx, orderID, provider)
}

func (s *Service) Webhook(ctx context.Context, provider, providerRef, status string, raw []byte) (*WebhookResult, error) {
	return s.repo.HandleWebhook(ctx, provider, providerRef, status, raw)
}
