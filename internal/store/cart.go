package store

import (
	"context"

	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
	"github.com/jackc/pgx/v5"
)

// GetCart returns the current cart for a user as a list of cart items.
func (s *Store) GetCart(ctx context.Context, userID int64) ([]types.CartItem, error) {
	query := `
		SELECT p.id AS id, p.id AS product_id, p.name, p.price, p.image_url, ci.qty
		FROM carts c
		JOIN cart_items ci ON ci.cart_id = c.id
		JOIN products p ON p.id = ci.product_id
		WHERE c.user_id = $1
		ORDER BY ci.id ASC
	`

	rows, err := s.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]types.CartItem, 0)
	for rows.Next() {
		var item types.CartItem
		if err := rows.Scan(&item.ID, &item.ProductID, &item.Name, &item.Price, &item.Image, &item.Qty); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// ReplaceCart overwrites the user's cart with the provided set of items.
func (s *Store) ReplaceCart(ctx context.Context, userID int64, items []types.CartItem) ([]types.CartItem, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	cartID, err := s.getOrCreateCartIDTx(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM cart_items WHERE cart_id = $1`, cartID); err != nil {
		return nil, err
	}

	for _, item := range items {
		productID := item.ProductID
		if productID == 0 {
			productID = item.ID
		}
		if productID == 0 || item.Qty <= 0 {
			continue
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO cart_items (cart_id, product_id, qty)
			VALUES ($1, $2, $3)
		`, cartID, productID, item.Qty); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.GetCart(ctx, userID)
}

func (s *Store) getOrCreateCartIDTx(ctx context.Context, tx pgx.Tx, userID int64) (int64, error) {
	var cartID int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO carts (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING
		RETURNING id
	`, userID).Scan(&cartID); err != nil {
		if err == pgx.ErrNoRows {
			if err := tx.QueryRow(ctx, `SELECT id FROM carts WHERE user_id = $1`, userID).Scan(&cartID); err != nil {
				return 0, err
			}
			return cartID, nil
		}
		return 0, err
	}
	return cartID, nil
}
