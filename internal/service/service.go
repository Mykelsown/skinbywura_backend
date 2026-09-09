// Package service contains the business logic layer for the application.
// Service functions coordinate validation, persistence, and API workflows.
package service

import "github.com/Mykelsown/skinbywura_backend.git/internal/store"

type Service struct {
	store *store.Store
}

func NewService(s *store.Store) (*Service) {
	return &Service{store: s}
}