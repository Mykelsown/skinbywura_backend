package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Mykelsown/skinbywura_backend.git/internal/service"
	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

// AuthHandler exposes authentication-related HTTP handlers.
type AuthHandler struct {
	service *service.Service
}

// NewAuthHandler creates an auth handler from a service dependency.
func NewAuthHandler(s *service.Service) *AuthHandler {
	return &AuthHandler{service: s}
}

// Signup creates a user and sets a session cookie on success.
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req types.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, token, err := h.service.Signup(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			WriteErrorJSON(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrEmailInUse):
			WriteErrorJSON(w, http.StatusConflict, err.Error())
		default:
			WriteErrorJSON(w, http.StatusInternalServerError, "failed to create account")
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     service.SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	WriteJSON(w, http.StatusCreated, user)
}

// Login authenticates a user and issues a session cookie.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, token, err := h.service.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			WriteErrorJSON(w, http.StatusUnauthorized, err.Error())
		default:
			WriteErrorJSON(w, http.StatusInternalServerError, "failed to login")
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     service.SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	WriteJSON(w, http.StatusOK, user)
}

// Logout clears any current session cookie.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(service.SessionCookieName)
	if err == nil && cookie.Value != "" {
		if err := h.service.Logout(r.Context(), cookie.Value); err != nil {
			WriteErrorJSON(w, http.StatusInternalServerError, "failed to logout")
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     service.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
	})

	WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

// Me returns the current authenticated user from middleware context.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(types.UserContextKey).(types.User)
	if !ok {
		WriteErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	WriteJSON(w, http.StatusOK, user)
}
