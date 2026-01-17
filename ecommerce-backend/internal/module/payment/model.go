package payment

type InitiateResult struct {
	PaymentID   string `json:"payment_id"`
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	Amount      int64  `json:"amount"`
	Provider    string `json:"provider"`
	ProviderRef string `json:"provider_ref"`
	PayURL      string `json:"pay_url,omitempty"`
}

type WebhookResult struct {
	PaymentID string `json:"payment_id"`
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
}
