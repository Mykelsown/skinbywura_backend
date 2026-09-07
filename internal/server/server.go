// Package server sets up the HTTP server instance and the routing layer for the app.
package server

import (
	"net/http"

	"github.com/Mykelsown/skinbywura_backend.git/config"
)

// Server holds the HTTP listener address and the registered handlers.
type Server struct {
	address string
	route   http.Handler
}

// New creates a server instance from the loaded environment configuration.
func New(cfg config.EnvData) *Server {
	ser := &Server{
		address: cfg.Port,
		route:   loadRoute(),
	}

	return ser
}

// Run starts the HTTP server on the configured port.
func (s *Server) Run() error {
	err := http.ListenAndServe(":"+s.address, s.route)
	return err
}
