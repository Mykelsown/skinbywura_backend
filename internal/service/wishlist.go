package service

import (
	"context"

	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

// GetWishlist loads a user's wishlist from storage.
func (s *Service) GetWishlist(ctx context.Context, userID int64) ([]types.WishlistItem, error) {
	return s.store.GetWishlist(ctx, userID)
}

// ReplaceWishlist saves the entire wishlist state for a user.
func (s *Service) ReplaceWishlist(ctx context.Context, userID int64, items []types.WishlistItem) ([]types.WishlistItem, error) {
	return s.store.ReplaceWishlist(ctx, userID, items)
}
