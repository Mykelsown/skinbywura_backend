package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/Mykelsown/skinbywura_backend.git/internal/service"
	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

// CookieSettings controls how the session cookie is issued.
type CookieSettings struct {
	Secure   bool
	SameSite http.SameSite
}

// ParseSameSite converts the config string into the net/http value.
func ParseSameSite(v string) http.SameSite {
	switch v {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

// AuthHandler exposes authentication-related HTTP handlers.
type AuthHandler struct {
	service *service.Service
	cookie  CookieSettings
}

// NewAuthHandler creates an auth handler from a service dependency and cookie policy.
func NewAuthHandler(s *service.Service, c CookieSettings) *AuthHandler {
	return &AuthHandler{service: s, cookie: c}
}

func (h *AuthHandler) setSession(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     service.SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookie.Secure,
		SameSite: h.cookie.SameSite,
		Expires:  expires,
		MaxAge:   maxAge(expires),
	})
}

func maxAge(expires time.Time) int {
	if expires.Unix() <= 0 {
		return -1
	}
	return int(time.Until(expires).Seconds())
}

// Signup creates a user and sets a session cookie on success.
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req types.SignupRequest
	if err := decodeJSON(w, r, &req); err != nil {
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
			serverError(w, r, "failed to create account", err)
		}
		return
	}

	h.setSession(w, token, service.SessionExpiry())
	WriteJSON(w, http.StatusCreated, user)
}

// Login authenticates a user and issues a session cookie.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := decodeJSON(w, r, &req); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, token, err := h.service.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			WriteErrorJSON(w, http.StatusUnauthorized, err.Error())
		default:
			serverError(w, r, "failed to login", err)
		}
		return
	}

	h.setSession(w, token, service.SessionExpiry())
	WriteJSON(w, http.StatusOK, user)
}

// Logout clears any current session cookie.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(service.SessionCookieName)
	if err == nil && cookie.Value != "" {
		if err := h.service.Logout(r.Context(), cookie.Value); err != nil {
			serverError(w, r, "failed to logout", err)
			return
		}
	}

	h.setSession(w, "", time.Unix(0, 0))
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

// UpdateMe changes the display name of the signed-in user.
func (h *AuthHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(types.UserContextKey).(types.User)
	if !ok {
		WriteErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req types.UpdateProfileRequest
	if err := decodeJSON(w, r, &req); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}
	updated, err := h.service.UpdateProfile(r.Context(), user, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			WriteErrorJSON(w, http.StatusBadRequest, "name is required")
			return
		}
		serverError(w, r, "failed to update profile", err)
		return
	}
	WriteJSON(w, http.StatusOK, updated)
}

// ChangePassword updates the password after checking the current one.
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(types.UserContextKey).(types.User)
	if !ok {
		WriteErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req types.ChangePasswordRequest
	if err := decodeJSON(w, r, &req); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}
	cookie, err := r.Cookie(service.SessionCookieName)
	if err != nil {
		WriteErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.service.ChangePassword(r.Context(), user, cookie.Value, req); err != nil {
		switch {
		case errors.Is(err, service.ErrWeakPassword):
			WriteErrorJSON(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrWrongPassword):
			WriteErrorJSON(w, http.StatusForbidden, err.Error())
		default:
			serverError(w, r, "failed to change password", err)
		}
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"message": "password updated"})
}

// ForgotPassword emails a reset link. The response is identical whether or not the account exists.
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req types.ForgotPasswordRequest
	if err := decodeJSON(w, r, &req); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.service.RequestPasswordReset(r.Context(), req); err != nil {
		serverError(w, r, "failed to process request", err)
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"message": "if that account exists, a reset link has been sent"})
}

// ResetPassword completes the reset flow with the emailed token.
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req types.ResetPasswordRequest
	if err := decodeJSON(w, r, &req); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.service.ResetPassword(r.Context(), req); err != nil {
		switch {
		case errors.Is(err, service.ErrWeakPassword), errors.Is(err, service.ErrResetInvalid):
			WriteErrorJSON(w, http.StatusBadRequest, err.Error())
		default:
			serverError(w, r, "failed to reset password", err)
		}
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"message": "password updated, please sign in"})
}
