// Package server contains the route registration and HTTP handlers for the API.
package server

import (
	"net/http"
	"time"

	"github.com/Mykelsown/skinbywura_backend.git/internal/handler"
	"github.com/Mykelsown/skinbywura_backend.git/internal/middleware"
	"github.com/Mykelsown/skinbywura_backend.git/internal/service"
)

// handlers groups every HTTP handler so route registration has one parameter.
type handlers struct {
	product  *handler.ProductHandler
	auth     *handler.AuthHandler
	cart     *handler.CartHandler
	wishlist *handler.WishlistHandler
	order    *handler.OrderHandler
	contact  *handler.ContactHandler
}

// loadRoute registers all API endpoints and returns the mux that serves them.
func loadRoute(svc *service.Service, h handlers) http.Handler {
	mux := http.NewServeMux()

	authed := func(fn http.HandlerFunc) http.Handler {
		return middleware.RequireAuth(svc)(fn)
	}

	// Credential and spam-prone endpoints get a per-IP quota.
	strict := middleware.NewRateLimiter(10, time.Minute)
	limited := func(fn http.HandlerFunc) http.Handler { return strict.Wrap(fn) }

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		handler.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /api/products", h.product.ListProducts)
	mux.HandleFunc("GET /api/products/{id}", h.product.GetProductByID)

	mux.Handle("POST /api/auth/signup", limited(h.auth.Signup))
	mux.Handle("POST /api/auth/login", limited(h.auth.Login))
	mux.Handle("POST /api/auth/forgot-password", limited(h.auth.ForgotPassword))
	mux.Handle("POST /api/auth/reset-password", limited(h.auth.ResetPassword))
	mux.Handle("POST /api/auth/logout", authed(h.auth.Logout))
	mux.Handle("GET /api/auth/me", authed(h.auth.Me))
	mux.Handle("PATCH /api/auth/me", authed(h.auth.UpdateMe))
	mux.Handle("POST /api/auth/change-password", limited(authed(h.auth.ChangePassword).ServeHTTP))

	mux.Handle("POST /api/contact", limited(h.contact.Submit))

	mux.Handle("GET /api/cart", authed(h.cart.GetCart))
	mux.Handle("PUT /api/cart", authed(h.cart.ReplaceCart))
	mux.Handle("GET /api/wishlist", authed(h.wishlist.GetWishlist))
	mux.Handle("PUT /api/wishlist", authed(h.wishlist.ReplaceWishlist))
	mux.Handle("POST /api/orders", authed(h.order.Checkout))
	mux.Handle("GET /api/orders", authed(h.order.ListOrders))

	return mux
}
