package payment

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")
var ErrInvalidStatus = errors.New("invalid status")

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) InitiatePayment(ctx context.Context, orderID, provider string) (*InitiateResult, error) {
	// get order total + status
	var status string
	var amount int64
	err := r.pool.QueryRow(ctx, `
SELECT status, grand_total
FROM orders
WHERE id=$1
LIMIT 1;
`, orderID).Scan(&status, &amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// only allow initiating for pending_payment (simple rule for now)
	if status != "pending_payment" {
		return nil, ErrNotFound
	}

	providerRef := newRef()

	// create payment row
	var paymentID string
	var payURL *string
	// for provider "manual", no pay_url; later provider can return URL
	err = r.pool.QueryRow(ctx, `
INSERT INTO payments (order_id, provider, status, amount, provider_ref, pay_url)
VALUES ($1, $2, 'initiated', $3, $4, $5)
RETURNING id::text;
`, orderID, provider, amount, providerRef, payURL).Scan(&paymentID)
	if err != nil {
		return nil, err
	}

	return &InitiateResult{
		PaymentID:   paymentID,
		OrderID:     orderID,
		Status:      "initiated",
		Amount:      amount,
		Provider:    provider,
		ProviderRef: providerRef,
	}, nil
}

func (r *PostgresRepository) HandleWebhook(ctx context.Context, provider, providerRef, newStatus string, rawPayload []byte) (*WebhookResult, error) {
	// validate status
	switch newStatus {
	case "pending", "paid", "failed", "expired", "refunded":
	default:
		return nil, ErrInvalidStatus
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// lock payment row
	var paymentID, orderID, curStatus string
	err = tx.QueryRow(ctx, `
SELECT id::text, order_id::text, status
FROM payments
WHERE provider=$1 AND provider_ref=$2
FOR UPDATE;
`, provider, providerRef).Scan(&paymentID, &orderID, &curStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// store payload
	var payloadJSON any
	_ = json.Unmarshal(rawPayload, &payloadJSON)
	_, err = tx.Exec(ctx, `
UPDATE payments
SET status=$1, payload=$2, updated_at=now()
WHERE provider=$3 AND provider_ref=$4;
`, newStatus, payloadJSON, provider, providerRef)
	if err != nil {
		return nil, err
	}

	// map payment status -> order status (simple)
	if newStatus == "paid" {
		_, err = tx.Exec(ctx, `
UPDATE orders SET status='paid', updated_at=now()
WHERE id=$1 AND status='pending_payment';
`, orderID)
		if err != nil {
			return nil, err
		}
	}
	if newStatus == "failed" || newStatus == "expired" {
		_, err = tx.Exec(ctx, `
UPDATE orders SET status='canceled', updated_at=now()
WHERE id=$1 AND status='pending_payment';
`, orderID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &WebhookResult{
		PaymentID: paymentID,
		OrderID:   orderID,
		Status:    newStatus,
	}, nil
}

func newRef() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
