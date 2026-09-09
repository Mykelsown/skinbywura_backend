// Package server contains the route registration and HTTP handlers for the API.
package server

import (
	"net/http"

	"github.com/Mykelsown/skinbywura_backend.git/internal/handler"
)

// loadRoute registers all API endpoints and returns the mux that serves them.
func loadRoute(productHandler *handler.ProductHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		handler.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /api/products", productHandler.ListProducts)
	mux.HandleFunc("GET /api/products/{id}", productHandler.GetProductByID)

	return mux
}
