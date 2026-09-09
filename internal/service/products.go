package service

import (
	"context"

	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

// ListProducts delegates to the store layer for now.
func (s *Service) ListProducts(ctx context.Context) ([]types.Product, error) {
	return s.store.ListProducts(ctx)
}

// GetProductByID delegates to the store layer for now.
func (s *Service) GetProductByID(ctx context.Context, id int) (types.Product, error) {
	return s.store.GetProductByID(ctx, id)
}
