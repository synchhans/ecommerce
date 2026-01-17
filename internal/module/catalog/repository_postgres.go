package catalog

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) ListProducts(ctx context.Context, limit, offset int, search string) ([]ProductListItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// simple search: match name/slug
	where := "p.is_active = true"
	args := []any{}
	argn := 1
	if strings.TrimSpace(search) != "" {
		where += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.slug ILIKE $%d)", argn, argn)
		args = append(args, "%"+strings.TrimSpace(search)+"%")
		argn++
	}

	q := fmt.Sprintf(`
SELECT
  p.id::text,
  p.slug,
  p.name,
  COALESCE(MIN(v.price), 0) AS min_price,
  COALESCE(MAX(v.price), 0) AS max_price,
  COALESCE((
     SELECT url FROM product_images pi
     WHERE pi.product_id = p.id
     ORDER BY pi.position ASC
     LIMIT 1
  ), '') AS image_url
FROM products p
LEFT JOIN product_variants v ON v.product_id = p.id AND v.is_active = true
WHERE %s
GROUP BY p.id, p.slug, p.name
ORDER BY p.created_at DESC
LIMIT $%d OFFSET $%d;
`, where, argn, argn+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]ProductListItem, 0, limit)
	for rows.Next() {
		var it ProductListItem
		if err := rows.Scan(&it.ID, &it.Slug, &it.Name, &it.MinPrice, &it.MaxPrice, &it.ImageURL); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *PostgresRepository) GetProductBySlug(ctx context.Context, slug string) (*ProductDetail, error) {
	// Product
	var p ProductDetail
	err := r.pool.QueryRow(ctx, `
SELECT id::text, slug, name, description
FROM products
WHERE slug = $1 AND is_active = true
LIMIT 1;
`, slug).Scan(&p.ID, &p.Slug, &p.Name, &p.Description)
	if err != nil {
		return nil, err
	}

	// Images
	imgRows, err := r.pool.Query(ctx, `
SELECT url, position
FROM product_images
WHERE product_id = $1
ORDER BY position ASC;
`, p.ID)
	if err != nil {
		return nil, err
	}
	defer imgRows.Close()

	for imgRows.Next() {
		var im ProductImage
		if err := imgRows.Scan(&im.URL, &im.Position); err != nil {
			return nil, err
		}
		p.Images = append(p.Images, im)
	}
	if err := imgRows.Err(); err != nil {
		return nil, err
	}

	// Variants
	varRows, err := r.pool.Query(ctx, `
SELECT id::text, sku, name, price, is_active
FROM product_variants
WHERE product_id = $1
ORDER BY created_at ASC;
`, p.ID)
	if err != nil {
		return nil, err
	}
	defer varRows.Close()

	for varRows.Next() {
		var v ProductVariant
		if err := varRows.Scan(&v.ID, &v.SKU, &v.Name, &v.Price, &v.Active); err != nil {
			return nil, err
		}
		p.Variants = append(p.Variants, v)
	}
	if err := varRows.Err(); err != nil {
		return nil, err
	}

	return &p, nil
}
