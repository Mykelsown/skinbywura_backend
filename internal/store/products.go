package store

import (
	"context"

	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

// ListProducts returns every product in the database in a stable order. it essentaially connects the DB to the stuct for the products built with GO.
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
