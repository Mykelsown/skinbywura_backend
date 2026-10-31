package middleware

import (
	"context"
	"net/http"

	"github.com/Mykelsown/skinbywura_backend.git/internal/service"
	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

// RequireAuth ensures the request carries a valid session cookie and attaches the user to context.
func RequireAuth(svc *service.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(service.SessionCookieName)
			if err != nil {
				unauthorized(w)
				return
			}

			user, err := svc.ValidateSession(r.Context(), cookie.Value)
			if err != nil {
				unauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), types.UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// unauthorized answers in the same JSON shape every other error uses, so the frontend can parse it.
func unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
}
