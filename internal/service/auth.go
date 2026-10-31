package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Mykelsown/skinbywura_backend.git/internal/store"
	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidInput       = errors.New("name, a valid email and a password of 8 to 72 characters are required")
	ErrWeakPassword       = errors.New("password must be between 8 and 72 characters")
	ErrWrongPassword      = errors.New("current password is incorrect")
	ErrResetInvalid       = errors.New("this reset link is invalid or has expired")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailInUse         = errors.New("email already in use")
	ErrUnauthorized       = errors.New("unauthorized")
)

const SessionCookieName = "session_token"

const sessionTTL = 7 * 24 * time.Hour

// SessionExpiry returns the expiry time for a session created now.
func SessionExpiry() time.Time { return time.Now().Add(sessionTTL) }

// Signup creates a user, hashes the password, and issues a session token.
func (s *Service) Signup(ctx context.Context, req types.SignupRequest) (types.User, string, error) {
	name := clean(req.Name, 80)
	email := strings.TrimSpace(strings.ToLower(req.Email))
	password := req.Password
	if name == "" || !validEmail(email) || !validPassword(password) {
		return types.User{}, "", ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return types.User{}, "", err
	}

	user, err := s.store.CreateUser(ctx, name, email, string(hash))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return types.User{}, "", ErrEmailInUse
		}
		return types.User{}, "", err
	}

	token, err := s.store.CreateSession(ctx, user.ID, time.Now().Add(sessionTTL))
	if err != nil {
		return types.User{}, "", err
	}

	return user, token, nil
}

// Login validates credentials and issues a session token when they are valid.
func (s *Service) Login(ctx context.Context, req types.LoginRequest) (types.User, string, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	password := req.Password
	if email == "" || password == "" {
		return types.User{}, "", ErrInvalidCredentials
	}

	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.User{}, "", ErrInvalidCredentials
		}
		return types.User{}, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return types.User{}, "", ErrInvalidCredentials
	}

	token, err := s.store.CreateSession(ctx, user.ID, time.Now().Add(sessionTTL))
	if err != nil {
		return types.User{}, "", err
	}

	return user, token, nil
}

// Logout removes a session token from storage.
func (s *Service) Logout(ctx context.Context, token string) error {
	return s.store.DeleteSession(ctx, token)
}

// ValidateSession confirms the token belongs to a current user and returns that user.
func (s *Service) ValidateSession(ctx context.Context, token string) (types.User, error) {
	if strings.TrimSpace(token) == "" {
		return types.User{}, ErrUnauthorized
	}

	userID, err := s.store.GetSession(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, store.ErrSessionExpired) {
			return types.User{}, ErrUnauthorized
		}
		return types.User{}, err
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.User{}, ErrUnauthorized
		}
		return types.User{}, err
	}

	return user, nil
}
