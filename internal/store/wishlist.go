package store

import (
	"context"

	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

// GetWishlist returns the products saved by a user.
func (s *Store) GetWishlist(ctx context.Context, userID int64) ([]types.WishlistItem, error) {
	query := `
		SELECT p.id AS id, p.id AS product_id, p.name, p.price, p.image_url
		FROM wishlist_items wi
		JOIN products p ON p.id = wi.product_id
		WHERE wi.user_id = $1
		ORDER BY wi.created_at DESC
	`

	rows, err := s.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]types.WishlistItem, 0)
	for rows.Next() {
		var item types.WishlistItem
		if err := rows.Scan(&item.ID, &item.ProductID, &item.Name, &item.Price, &item.Image); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// ReplaceWishlist overwrites the user's wishlist with the provided items.
func (s *Store) ReplaceWishlist(ctx context.Context, userID int64, items []types.WishlistItem) ([]types.WishlistItem, error) {
	if _, err := s.Pool.Exec(ctx, `DELETE FROM wishlist_items WHERE user_id = $1`, userID); err != nil {
		return nil, err
	}

	for _, item := range items {
		productID := item.ProductID
		if productID == 0 {
			productID = item.ID
		}
		if productID == 0 {
			continue
		}
		if _, err := s.Pool.Exec(ctx, `
			INSERT INTO wishlist_items (user_id, product_id)
			VALUES ($1, $2)
			ON CONFLICT (user_id, product_id) DO NOTHING
		`, userID, productID); err != nil {
			return nil, err
		}
	}

	return s.GetWishlist(ctx, userID)
}
