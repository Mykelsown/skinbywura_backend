package handler

import (
	"errors"
	"net/http"

	"github.com/Mykelsown/skinbywura_backend.git/internal/service"
	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

// ContactHandler receives messages from the public contact form.
type ContactHandler struct {
	service *service.Service
}

// NewContactHandler creates a contact handler from a service dependency.
func NewContactHandler(s *service.Service) *ContactHandler {
	return &ContactHandler{service: s}
}

// Submit stores the message and alerts the store owner.
func (h *ContactHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var req types.ContactRequest
	if err := decodeJSON(w, r, &req); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.service.SubmitContact(r.Context(), req); err != nil {
		if errors.Is(err, service.ErrInvalidContact) {
			WriteErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		serverError(w, r, "failed to send message", err)
		return
	}
	WriteJSON(w, http.StatusCreated, map[string]string{"message": "message received"})
}
