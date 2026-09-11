package service

import (
	"context"

	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

// GetCart loads a user's cart from storage.
func (s *Service) GetCart(ctx context.Context, userID int64) ([]types.CartItem, error) {
	return s.store.GetCart(ctx, userID)
}

// ReplaceCart saves the full cart state for a user.
func (s *Service) ReplaceCart(ctx context.Context, userID int64, items []types.CartItem) ([]types.CartItem, error) {
	return s.store.ReplaceCart(ctx, userID, items)
}
