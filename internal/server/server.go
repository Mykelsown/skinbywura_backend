// Package server sets up the HTTP server instance and the routing layer for the app.
package server

import (
	"net/http"

	"github.com/Mykelsown/skinbywura_backend.git/config"
	"github.com/Mykelsown/skinbywura_backend.git/internal/handler"
	"github.com/Mykelsown/skinbywura_backend.git/internal/middleware"
	"github.com/Mykelsown/skinbywura_backend.git/internal/service"
)

// Server holds the HTTP listener address and the registered handlers.
type Server struct {
	address string       // address is the configured TCP port.
	route   http.Handler // route handles incoming HTTP requests.
}

// New creates a server instance from the loaded environment configuration.
func New(cfg config.EnvData, svc *service.Service) *Server {
	productHandler := handler.NewProductHandler(svc)

	ser := &Server{
		address: cfg.Port,
		route:   middleware.CORS(loadRoute(productHandler)),
	}

	return ser
}

// Run starts the HTTP server on the configured port.
func (s *Server) Run() error {
	// ListenAndServe blocks until the server stops or encounters an error.
	err := http.ListenAndServe(":"+s.address, s.route)
	return err
}
