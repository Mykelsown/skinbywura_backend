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

// UpdateProfileRequest changes the display name on an account.
type UpdateProfileRequest struct {
	Name string `json:"name"`
}

// ChangePasswordRequest changes the password for a signed-in user.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ForgotPasswordRequest asks for a reset link to be emailed.
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest completes a reset using the emailed token.
type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// ContactRequest is a message submitted from the public contact form.
type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
}
