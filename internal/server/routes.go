// Package server contains the route registration and HTTP handlers for the API.
package server

import (
	"net/http"

	"github.com/Mykelsown/skinbywura_backend.git/internal/handler"
	"github.com/Mykelsown/skinbywura_backend.git/internal/middleware"
	"github.com/Mykelsown/skinbywura_backend.git/internal/service"
)

// loadRoute registers all API endpoints and returns the mux that serves them.
func loadRoute(svc *service.Service, productHandler *handler.ProductHandler, authHandler *handler.AuthHandler, cartHandler *handler.CartHandler, wishlistHandler *handler.WishlistHandler, orderHandler *handler.OrderHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		handler.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /api/products", productHandler.ListProducts)
	mux.HandleFunc("GET /api/products/{id}", productHandler.GetProductByID)
	mux.HandleFunc("POST /api/auth/signup", authHandler.Signup)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.Handle("POST /api/auth/logout", middleware.RequireAuth(svc)(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("GET /api/auth/me", middleware.RequireAuth(svc)(http.HandlerFunc(authHandler.Me)))

	mux.Handle("GET /api/cart", middleware.RequireAuth(svc)(http.HandlerFunc(cartHandler.GetCart)))
	mux.Handle("PUT /api/cart", middleware.RequireAuth(svc)(http.HandlerFunc(cartHandler.ReplaceCart)))
	mux.Handle("GET /api/wishlist", middleware.RequireAuth(svc)(http.HandlerFunc(wishlistHandler.GetWishlist)))
	mux.Handle("PUT /api/wishlist", middleware.RequireAuth(svc)(http.HandlerFunc(wishlistHandler.ReplaceWishlist)))
	mux.Handle("POST /api/orders", middleware.RequireAuth(svc)(http.HandlerFunc(orderHandler.Checkout)))
	mux.Handle("GET /api/orders", middleware.RequireAuth(svc)(http.HandlerFunc(orderHandler.ListOrders)))

	return mux
}
