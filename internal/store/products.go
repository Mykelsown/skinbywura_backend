package store

import (
	"context"
	"errors"

	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
	"github.com/jackc/pgx/v5"
)

// ListProducts returns every product in the database in a stable order.
func (s *Store) ListProducts(ctx context.Context) ([]types.Product, error) {
	query := `
		SELECT
			id,
			name,
			category,
			price,
			COALESCE(compare_at_price, 0) AS compare_at_price,
			rating,
			review_count,
			skin_type,
			COALESCE(badge, '') AS badge,
			image_url,
			description,
			volume
		FROM products
		ORDER BY id
	`

	rows, err := s.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]types.Product, 0)
	for rows.Next() {
		var product types.Product
		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Category,
			&product.Price,
			&product.CompareAtPrice,
			&product.Rating,
			&product.ReviewCount,
			&product.SkinType,
			&product.Badge,
			&product.ImageURL,
			&product.Description,
			&product.Volume,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// GetProductByID returns a single product by its ID.
func (s *Store) GetProductByID(ctx context.Context, id int) (types.Product, error) {
	query := `
		SELECT
			id,
			name,
			category,
			price,
			COALESCE(compare_at_price, 0) AS compare_at_price,
			rating,
			review_count,
			skin_type,
			COALESCE(badge, '') AS badge,
			image_url,
			description,
			volume
		FROM products
		WHERE id = $1
	`

	var product types.Product
	err := s.Pool.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Category,
		&product.Price,
		&product.CompareAtPrice,
		&product.Rating,
		&product.ReviewCount,
		&product.SkinType,
		&product.Badge,
		&product.ImageURL,
		&product.Description,
		&product.Volume,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.Product{}, pgx.ErrNoRows
		}
		return types.Product{}, err
	}

	return product, nil
}
