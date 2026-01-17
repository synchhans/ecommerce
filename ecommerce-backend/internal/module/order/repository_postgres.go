package order

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")
var ErrEmptyCart = errors.New("empty cart")

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateOrderFromCart(ctx context.Context, cartID string, shipAddr AddressSnapshot) (string, error) {
	// Transaction is important.
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// 1) lock cart row (simple)
	var cartStatus string
	err = tx.QueryRow(ctx, `SELECT status FROM carts WHERE id=$1 FOR UPDATE`, cartID).Scan(&cartStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	if cartStatus != "active" {
		// treat as not found / invalid for now
		return "", ErrNotFound
	}

	// 2) get cart items
	rows, err := tx.Query(ctx, `
SELECT ci.variant_id::text, ci.qty, v.sku, v.name, v.price
FROM cart_items ci
JOIN product_variants v ON v.id = ci.variant_id
JOIN products p ON p.id = v.product_id
WHERE ci.cart_id = $1 AND v.is_active = true AND p.is_active = true
ORDER BY ci.created_at ASC;
`, cartID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	type itemRow struct {
		variantID string
		qty       int
		sku       string
		name      string
		price     int64
	}
	var items []itemRow
	var subtotal int64

	for rows.Next() {
		var it itemRow
		if err := rows.Scan(&it.variantID, &it.qty, &it.sku, &it.name, &it.price); err != nil {
			return "", err
		}
		if it.qty <= 0 {
			continue
		}
		line := it.price * int64(it.qty)
		subtotal += line
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "", ErrEmptyCart
	}

	discountTotal := int64(0)
	shippingTotal := int64(0)
	grandTotal := subtotal - discountTotal + shippingTotal

	orderNumber := generateOrderNumber()

	addrJSON, _ := json.Marshal(shipAddr)

	// 3) create order
	var orderID string
	err = tx.QueryRow(ctx, `
INSERT INTO orders (order_number, cart_id, status, currency, subtotal, discount_total, shipping_total, grand_total, shipping_address_snapshot)
VALUES ($1, $2, 'pending_payment', 'IDR', $3, $4, $5, $6, $7)
RETURNING id::text;
`, orderNumber, cartID, subtotal, discountTotal, shippingTotal, grandTotal, addrJSON).Scan(&orderID)
	if err != nil {
		return "", err
	}

	// 4) create order items
	for _, it := range items {
		lineTotal := it.price * int64(it.qty)
		_, err := tx.Exec(ctx, `
INSERT INTO order_items (order_id, variant_id, sku, name, unit_price, qty, line_total)
VALUES ($1, $2, $3, $4, $5, $6, $7);
`, orderID, it.variantID, it.sku, it.name, it.price, it.qty, lineTotal)
		if err != nil {
			return "", err
		}
	}

	// 5) mark cart converted (optional but useful)
	_, err = tx.Exec(ctx, `UPDATE carts SET status='converted', updated_at=now() WHERE id=$1;`, cartID)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}
	return orderID, nil
}

func (r *PostgresRepository) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	var o Order
	err := r.pool.QueryRow(ctx, `
SELECT id::text, order_number, status, currency, subtotal, discount_total, shipping_total, grand_total
FROM orders
WHERE id=$1
LIMIT 1;
`, orderID).Scan(&o.ID, &o.OrderNumber, &o.Status, &o.Currency, &o.Subtotal, &o.Discount, &o.Shipping, &o.GrandTotal)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	rows, err := r.pool.Query(ctx, `
SELECT id::text, variant_id::text, sku, name, unit_price, qty, line_total
FROM order_items
WHERE order_id=$1
ORDER BY id ASC;
`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var it OrderItem
		if err := rows.Scan(&it.ID, &it.VariantID, &it.SKU, &it.Name, &it.UnitPrice, &it.Qty, &it.LineTotal); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &o, nil
}

// generateOrderNumber: human-friendly, unique enough for small-medium scale.
// Example: EC-20260112-8F3A2C
func generateOrderNumber() string {
	date := time.Now().Format("20060102")
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	return fmt.Sprintf("EC-%s-%X", date, b)
}
