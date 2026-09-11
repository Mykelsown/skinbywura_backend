package store

import (
	"context"
	"errors"

	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

var ErrCartEmpty = errors.New("cart is empty")

// CreateOrder creates a single order from the authenticated user's cart.
func (s *Store) CreateOrder(ctx context.Context, userID int64) (types.Order, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return types.Order{}, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		SELECT p.id, p.name, p.price, ci.qty
		FROM carts c
		JOIN cart_items ci ON ci.cart_id = c.id
		JOIN products p ON p.id = ci.product_id
		WHERE c.user_id = $1
		ORDER BY ci.id ASC
	`, userID)
	if err != nil {
		return types.Order{}, err
	}
	defer rows.Close()

	cartItems := make([]types.CartItem, 0)
	total := 0
	for rows.Next() {
		var item types.CartItem
		if err := rows.Scan(&item.ProductID, &item.Name, &item.Price, &item.Qty); err != nil {
			return types.Order{}, err
		}
		cartItems = append(cartItems, item)
		total += item.Price * item.Qty
	}
	if err := rows.Err(); err != nil {
		return types.Order{}, err
	}
	if len(cartItems) == 0 {
		return types.Order{}, ErrCartEmpty
	}

	var orderID int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO orders (user_id, status, total)
		VALUES ($1, 'pending', $2)
		RETURNING id
	`, userID, total).Scan(&orderID); err != nil {
		return types.Order{}, err
	}

	for _, item := range cartItems {
		if _, err := tx.Exec(ctx, `
			INSERT INTO order_items (order_id, product_id, product_name, product_price, qty)
			VALUES ($1, $2, $3, $4, $5)
		`, orderID, item.ProductID, item.Name, item.Price, item.Qty); err != nil {
			return types.Order{}, err
		}
	}

	if _, err := tx.Exec(ctx, `
		DELETE FROM cart_items
		WHERE cart_id = (SELECT id FROM carts WHERE user_id = $1)
	`, userID); err != nil {
		return types.Order{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return types.Order{}, err
	}

	return s.getOrderByID(ctx, orderID)
}

// ListOrdersByUser returns all orders for a user, newest first.
func (s *Store) ListOrdersByUser(ctx context.Context, userID int64) ([]types.Order, error) {
	rows, err := s.Pool.Query(ctx, `
		SELECT id, status, total, created_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC, id DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]types.Order, 0)
	for rows.Next() {
		var order types.Order
		if err := rows.Scan(&order.ID, &order.Status, &order.Total, &order.CreatedAt); err != nil {
			return nil, err
		}

		items, err := s.getOrderItemsByOrderID(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Store) getOrderByID(ctx context.Context, orderID int64) (types.Order, error) {
	var order types.Order
	if err := s.Pool.QueryRow(ctx, `
		SELECT id, status, total, created_at
		FROM orders
		WHERE id = $1
	`, orderID).Scan(&order.ID, &order.Status, &order.Total, &order.CreatedAt); err != nil {
		return types.Order{}, err
	}

	items, err := s.getOrderItemsByOrderID(ctx, orderID)
	if err != nil {
		return types.Order{}, err
	}
	order.Items = items
	return order, nil
}

func (s *Store) getOrderItemsByOrderID(ctx context.Context, orderID int64) ([]types.OrderItem, error) {
	rows, err := s.Pool.Query(ctx, `
		SELECT id, product_id, product_name, product_price, qty
		FROM order_items
		WHERE order_id = $1
		ORDER BY id ASC
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]types.OrderItem, 0)
	for rows.Next() {
		var item types.OrderItem
		if err := rows.Scan(&item.ID, &item.ProductID, &item.Name, &item.Price, &item.Qty); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
