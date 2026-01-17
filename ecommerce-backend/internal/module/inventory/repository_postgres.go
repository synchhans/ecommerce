package inventory

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) GetAvailability(ctx context.Context, variantID string) (*Availability, error) {
	var a Availability
	err := r.pool.QueryRow(ctx, `
SELECT variant_id::text, stock_on_hand, reserved
FROM inventory_items
WHERE variant_id = $1
LIMIT 1;
`, variantID).Scan(&a.VariantID, &a.StockOnHand, &a.Reserved)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	a.Available = a.StockOnHand - a.Reserved
	if a.Available < 0 {
		a.Available = 0
	}
	return &a, nil
}
