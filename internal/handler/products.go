package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mykelsown/skinbywura_backend.git/internal/service"
	"github.com/jackc/pgx/v5"
)

// ProductHandler exposes product-related HTTP handlers backed by a service layer.
type ProductHandler struct {
	service *service.Service
}

// NewProductHandler creates a product handler from a service dependency.
func NewProductHandler(s *service.Service) *ProductHandler {
	return &ProductHandler{service: s}
}

// ListProducts returns all products as JSON.
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		WriteErrorJSON(w, http.StatusInternalServerError, "failed to list products")
		return
	}

	WriteJSON(w, http.StatusOK, products)
}

// GetProductByID returns a single product by its route ID.
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "invalid product id")
		return
	}

	product, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			WriteErrorJSON(w, http.StatusNotFound, "product not found")
			return
		}

		WriteErrorJSON(w, http.StatusInternalServerError, "failed to fetch product")
		return
	}

	WriteJSON(w, http.StatusOK, product)
}
