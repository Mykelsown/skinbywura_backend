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
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := svc.ValidateSession(r.Context(), cookie.Value)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), types.UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
