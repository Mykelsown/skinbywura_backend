package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Mykelsown/skinbywura_backend.git/internal/service"
	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

// CartHandler exposes cart-related HTTP handlers for authenticated users.
type CartHandler struct {
	service *service.Service
}

// NewCartHandler creates a cart handler from a service dependency.
func NewCartHandler(s *service.Service) *CartHandler {
	return &CartHandler{service: s}
}

// GetCart returns the authenticated user's cart.
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(types.UserContextKey).(types.User)
	if !ok {
		WriteErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	cart, err := h.service.GetCart(r.Context(), user.ID)
	if err != nil {
		WriteErrorJSON(w, http.StatusInternalServerError, "failed to load cart")
		return
	}

	WriteJSON(w, http.StatusOK, cart)
}

// ReplaceCart overwrites the authenticated user's cart.
func (h *CartHandler) ReplaceCart(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(types.UserContextKey).(types.User)
	if !ok {
		WriteErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var items []types.CartItem
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "invalid cart payload")
		return
	}

	updated, err := h.service.ReplaceCart(r.Context(), user.ID, items)
	if err != nil {
		WriteErrorJSON(w, http.StatusInternalServerError, "failed to save cart")
		return
	}

	WriteJSON(w, http.StatusOK, updated)
}
