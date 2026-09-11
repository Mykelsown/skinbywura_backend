package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Mykelsown/skinbywura_backend.git/internal/service"
	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

// WishlistHandler exposes wishlist-related HTTP handlers for authenticated users.
type WishlistHandler struct {
	service *service.Service
}

// NewWishlistHandler creates a wishlist handler from a service dependency.
func NewWishlistHandler(s *service.Service) *WishlistHandler {
	return &WishlistHandler{service: s}
}

// GetWishlist returns the authenticated user's wishlist.
func (h *WishlistHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(types.UserContextKey).(types.User)
	if !ok {
		WriteErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	items, err := h.service.GetWishlist(r.Context(), user.ID)
	if err != nil {
		WriteErrorJSON(w, http.StatusInternalServerError, "failed to load wishlist")
		return
	}

	WriteJSON(w, http.StatusOK, items)
}

// ReplaceWishlist overwrites the authenticated user's wishlist.
func (h *WishlistHandler) ReplaceWishlist(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(types.UserContextKey).(types.User)
	if !ok {
		WriteErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var items []types.WishlistItem
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "invalid wishlist payload")
		return
	}

	updated, err := h.service.ReplaceWishlist(r.Context(), user.ID, items)
	if err != nil {
		WriteErrorJSON(w, http.StatusInternalServerError, "failed to save wishlist")
		return
	}

	WriteJSON(w, http.StatusOK, updated)
}
