// Package service contains the business logic layer for the application.
// Service functions coordinate validation, persistence, and API workflows.
package service

import "github.com/Mykelsown/skinbywura_backend.git/internal/store"

// Service coordinates business logic and delegates persistence to the store layer.
type Service struct {
	store *store.Store
}

// New creates a service instance from a store dependency.
func New(s *store.Store) *Service {
	return &Service{store: s}
}

// NewService creates a service instance from a store dependency.
func NewService(s *store.Store) *Service {
	return New(s)
}
