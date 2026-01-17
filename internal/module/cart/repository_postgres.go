package cart

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")
var ErrInvalidQty = errors.New("invalid qty")

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateCart(ctx context.Context) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx, `INSERT INTO carts DEFAULT VALUES RETURNING id::text;`).Scan(&id)
	return id, err
}

func (r *PostgresRepository) GetCart(ctx context.Context, cartID string) (*Cart, error) {
	var c Cart
	err := r.pool.QueryRow(ctx, `
SELECT id::text, status
FROM carts
WHERE id = $1
LIMIT 1;
`, cartID).Scan(&c.ID, &c.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	rows, err := r.pool.Query(ctx, `
SELECT id::text, variant_id::text, qty
FROM cart_items
WHERE cart_id = $1
ORDER BY created_at ASC;
`, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var it CartItem
		if err := rows.Scan(&it.ID, &it.VariantID, &it.Qty); err != nil {
			return nil, err
		}
		c.Items = append(c.Items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *PostgresRepository) UpsertItem(ctx context.Context, cartID, variantID string, qty int) error {
	if qty <= 0 {
		return ErrInvalidQty
	}

	// Upsert based on UNIQUE(cart_id, variant_id)
	_, err := r.pool.Exec(ctx, `
INSERT INTO cart_items (cart_id, variant_id, qty)
VALUES ($1, $2, $3)
ON CONFLICT (cart_id, variant_id)
DO UPDATE SET qty = EXCLUDED.qty, updated_at = now();
`, cartID, variantID, qty)
	return err
}

func (r *PostgresRepository) UpdateItemQty(ctx context.Context, cartID, itemID string, qty int) error {
	if qty <= 0 {
		return ErrInvalidQty
	}
	ct, err := r.pool.Exec(ctx, `
UPDATE cart_items
SET qty = $1, updated_at = now()
WHERE id = $2 AND cart_id = $3;
`, qty, itemID, cartID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) DeleteItem(ctx context.Context, cartID, itemID string) error {
	ct, err := r.pool.Exec(ctx, `
DELETE FROM cart_items
WHERE id = $1 AND cart_id = $2;
`, itemID, cartID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
