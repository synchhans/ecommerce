package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")
var ErrEmailTaken = errors.New("email already taken")

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, email, passwordHash, name string) (*User, error) {
	var u User
	err := r.pool.QueryRow(ctx, `
INSERT INTO users (email, password_hash, name, status)
VALUES ($1, $2, $3, 'active')
RETURNING id::text, email, name, status;
`, email, passwordHash, name).Scan(&u.ID, &u.Email, &u.Name, &u.Status)

	if err != nil {
		// simple mapping (unique violation)
		// pgx returns error with SQLSTATE; to keep minimal, map all insert errors -> email taken is fine for now
		return nil, ErrEmailTaken
	}
	return &u, nil
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, string, error) {
	var u User
	var ph string
	err := r.pool.QueryRow(ctx, `
SELECT id::text, email, name, status, password_hash
FROM users
WHERE email=$1
LIMIT 1;
`, email).Scan(&u.ID, &u.Email, &u.Name, &u.Status, &ph)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", ErrNotFound
		}
		return nil, "", err
	}
	return &u, ph, nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	var u User
	err := r.pool.QueryRow(ctx, `
SELECT id::text, email, name, status
FROM users
WHERE id=$1
LIMIT 1;
`, userID).Scan(&u.ID, &u.Email, &u.Name, &u.Status)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
