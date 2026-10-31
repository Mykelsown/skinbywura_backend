// Package service contains the business logic layer for the application.
// Service functions coordinate validation, persistence, and API workflows.
package service

import (
	"github.com/Mykelsown/skinbywura_backend.git/internal/mailer"
	"github.com/Mykelsown/skinbywura_backend.git/internal/store"
)

// Options carries the non-database dependencies and settings the service needs.
type Options struct {
	Mailer      mailer.Sender // Mailer delivers transactional email.
	FrontendURL string        // FrontendURL builds links inside emails.
	StoreEmail  string        // StoreEmail receives contact messages and order alerts.
	StoreName   string        // StoreName appears in email copy.
	PaymentNote string        // PaymentNote is appended to order confirmation emails.
}

// Service coordinates business logic and delegates persistence to the store layer.
type Service struct {
	store *store.Store
	opts  Options
}

// New creates a service instance from a store dependency and its options.
func New(s *store.Store, opts Options) *Service {
	if opts.Mailer == nil {
		opts.Mailer = mailer.LogSender{}
	}
	if opts.StoreName == "" {
		opts.StoreName = "Skin By Wura"
	}
	return &Service{store: s, opts: opts}
}
