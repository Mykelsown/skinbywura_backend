// Package types defines shared data models used across the backend.
package types

import "time"

// User represents a persisted application user.
// The password hash is stored internally but intentionally omitted from JSON output.
type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

// SignupRequest contains the payload used for creating a new account.
type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest contains the payload used to authenticate a user.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ContextKey identifies the request-scoped user value in context.
type ContextKey string

const UserContextKey ContextKey = "user"
