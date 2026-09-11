package service

import (
	"context"
	"errors"

	"github.com/Mykelsown/skinbywura_backend.git/internal/store"
	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

var ErrCartEmpty = store.ErrCartEmpty

// CreateOrder creates a user order from the current authenticated cart.
func (s *Service) CreateOrder(ctx context.Context, userID int64) (types.Order, error) {
	order, err := s.store.CreateOrder(ctx, userID)
	if err != nil {
		if errors.Is(err, store.ErrCartEmpty) {
			return types.Order{}, ErrCartEmpty
		}
		return types.Order{}, err
	}
	return order, nil
}

// ListOrdersByUser returns all orders for the given user.
func (s *Service) ListOrdersByUser(ctx context.Context, userID int64) ([]types.Order, error) {
	return s.store.ListOrdersByUser(ctx, userID)
}
