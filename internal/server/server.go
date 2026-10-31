// Package server sets up the HTTP server instance and the routing layer for the app.
package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Mykelsown/skinbywura_backend.git/config"
	"github.com/Mykelsown/skinbywura_backend.git/internal/handler"
	"github.com/Mykelsown/skinbywura_backend.git/internal/middleware"
	"github.com/Mykelsown/skinbywura_backend.git/internal/service"
)

// Server wraps the configured http.Server.
type Server struct {
	http *http.Server
}

// New creates a server instance from the loaded environment configuration.
func New(cfg config.EnvData, svc *service.Service) *Server {
	cookie := handler.CookieSettings{Secure: cfg.CookieSecure, SameSite: handler.ParseSameSite(cfg.CookieSameSite)}

	h := handlers{
		product:  handler.NewProductHandler(svc),
		auth:     handler.NewAuthHandler(svc, cookie),
		cart:     handler.NewCartHandler(svc),
		wishlist: handler.NewWishlistHandler(svc),
		order:    handler.NewOrderHandler(svc),
		contact:  handler.NewContactHandler(svc),
	}

	// Order matters: logging wraps everything so panics and CORS rejections are recorded.
	var root http.Handler = loadRoute(svc, h)
	root = middleware.CORS(cfg.AllowedOrigins)(root)
	root = middleware.SecurityHeaders(root)
	root = middleware.Logging(root)

	return &Server{http: &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           root,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       90 * time.Second,
		MaxHeaderBytes:    1 << 16,
	}}
}

// Run serves until ctx is cancelled (SIGTERM from the platform), then drains in-flight requests.
func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		slog.Info("listening", "addr", s.http.Addr)
		errCh <- s.http.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		return s.http.Shutdown(shutdownCtx)
	}
}
