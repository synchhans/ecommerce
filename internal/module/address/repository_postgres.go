package address

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")
var ErrInvalidPayload = errors.New("invalid payload")

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, userID string, a Address) (*Address, error) {
	if userID == "" || a.RecipientName == "" || a.Phone == "" || a.AddressLine1 == "" || a.City == "" || a.Province == "" || a.PostalCode == "" {
		return nil, ErrInvalidPayload
	}
	if a.Country == "" {
		a.Country = "ID"
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// If new address is default, unset previous default
	if a.IsDefault {
		_, err := tx.Exec(ctx, `UPDATE user_addresses SET is_default=false, updated_at=now() WHERE user_id=$1 AND is_default=true;`, userID)
		if err != nil {
			return nil, err
		}
	}

	var out Address
	err = tx.QueryRow(ctx, `
INSERT INTO user_addresses (
  user_id, label, recipient_name, phone, address_line1, address_line2, city, province, postal_code, country, is_default
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
RETURNING id::text, label, recipient_name, phone, address_line1, COALESCE(address_line2,''), city, province, postal_code, country, is_default;
`, userID, a.Label, a.RecipientName, a.Phone, a.AddressLine1, a.AddressLine2, a.City, a.Province, a.PostalCode, a.Country, a.IsDefault).
		Scan(&out.ID, &out.Label, &out.RecipientName, &out.Phone, &out.AddressLine1, &out.AddressLine2, &out.City, &out.Province, &out.PostalCode, &out.Country, &out.IsDefault)

	if err != nil {
		return nil, err
	}

	// If user has no default yet and this is not default, make it default automatically
	if !out.IsDefault {
		var hasDefault bool
		err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM user_addresses WHERE user_id=$1 AND is_default=true);`, userID).Scan(&hasDefault)
		if err != nil {
			return nil, err
		}
		if !hasDefault {
			_, err = tx.Exec(ctx, `UPDATE user_addresses SET is_default=true, updated_at=now() WHERE id=$1 AND user_id=$2;`, out.ID, userID)
			if err != nil {
				return nil, err
			}
			out.IsDefault = true
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *PostgresRepository) List(ctx context.Context, userID string) ([]Address, error) {
	rows, err := r.pool.Query(ctx, `
SELECT id::text, label, recipient_name, phone, address_line1, COALESCE(address_line2,''), city, province, postal_code, country, is_default
FROM user_addresses
WHERE user_id=$1
ORDER BY is_default DESC, created_at DESC;
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Address
	for rows.Next() {
		var a Address
		if err := rows.Scan(&a.ID, &a.Label, &a.RecipientName, &a.Phone, &a.AddressLine1, &a.AddressLine2, &a.City, &a.Province, &a.PostalCode, &a.Country, &a.IsDefault); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *PostgresRepository) Update(ctx context.Context, userID, addressID string, a Address) (*Address, error) {
	if userID == "" || addressID == "" {
		return nil, ErrInvalidPayload
	}
	if a.Country == "" {
		a.Country = "ID"
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if a.IsDefault {
		_, err := tx.Exec(ctx, `UPDATE user_addresses SET is_default=false, updated_at=now() WHERE user_id=$1 AND is_default=true;`, userID)
		if err != nil {
			return nil, err
		}
	}

	ct, err := tx.Exec(ctx, `
UPDATE user_addresses
SET label=$1, recipient_name=$2, phone=$3, address_line1=$4, address_line2=$5,
    city=$6, province=$7, postal_code=$8, country=$9, is_default=$10, updated_at=now()
WHERE id=$11 AND user_id=$12;
`, a.Label, a.RecipientName, a.Phone, a.AddressLine1, a.AddressLine2, a.City, a.Province, a.PostalCode, a.Country, a.IsDefault, addressID, userID)
	if err != nil {
		return nil, err
	}
	if ct.RowsAffected() == 0 {
		return nil, ErrNotFound
	}

	var out Address
	err = tx.QueryRow(ctx, `
SELECT id::text, label, recipient_name, phone, address_line1, COALESCE(address_line2,''), city, province, postal_code, country, is_default
FROM user_addresses
WHERE id=$1 AND user_id=$2
LIMIT 1;
`, addressID, userID).
		Scan(&out.ID, &out.Label, &out.RecipientName, &out.Phone, &out.AddressLine1, &out.AddressLine2, &out.City, &out.Province, &out.PostalCode, &out.Country, &out.IsDefault)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, userID, addressID string) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM user_addresses WHERE id=$1 AND user_id=$2;`, addressID, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) SetDefault(ctx context.Context, userID, addressID string) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `UPDATE user_addresses SET is_default=false, updated_at=now() WHERE user_id=$1 AND is_default=true;`, userID)
	if err != nil {
		return err
	}

	ct, err := tx.Exec(ctx, `
UPDATE user_addresses SET is_default=true, updated_at=now()
WHERE id=$1 AND user_id=$2;
`, addressID, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}

	return tx.Commit(ctx)
}
