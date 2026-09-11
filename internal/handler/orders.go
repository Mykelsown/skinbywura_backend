package handler

import (
	"errors"
	"net/http"

	"github.com/Mykelsown/skinbywura_backend.git/internal/service"
	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

// OrderHandler exposes order-related HTTP handlers for authenticated users.
type OrderHandler struct {
	service *service.Service
}

// NewOrderHandler creates an order handler from a service dependency.
func NewOrderHandler(s *service.Service) *OrderHandler {
	return &OrderHandler{service: s}
}

// Checkout creates an order from the authenticated user's current cart.
func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(types.UserContextKey).(types.User)
	if !ok {
		WriteErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	order, err := h.service.CreateOrder(r.Context(), user.ID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCartEmpty):
			WriteErrorJSON(w, http.StatusBadRequest, err.Error())
		default:
			WriteErrorJSON(w, http.StatusInternalServerError, "failed to create order")
		}
		return
	}

	WriteJSON(w, http.StatusCreated, order)
}

// ListOrders returns all orders for the authenticated user.
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(types.UserContextKey).(types.User)
	if !ok {
		WriteErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orders, err := h.service.ListOrdersByUser(r.Context(), user.ID)
	if err != nil {
		WriteErrorJSON(w, http.StatusInternalServerError, "failed to load orders")
		return
	}

	WriteJSON(w, http.StatusOK, orders)
}
