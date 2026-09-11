package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

var ErrSessionExpired = errors.New("session expired")

// generateToken creates a cryptographically random session token.
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreateSession inserts a new session row for the user and returns the token.
func (s *Store) CreateSession(ctx context.Context, userID int64, expiresAt time.Time) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	query := `
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var createdToken string
	err = s.Pool.QueryRow(ctx, query, token, userID, expiresAt).Scan(&createdToken)
	if err != nil {
		return "", err
	}

	return createdToken, nil
}

// GetSession returns the user ID for a valid, unexpired session token.
func (s *Store) GetSession(ctx context.Context, token string) (int64, error) {
	query := `
		SELECT user_id
		FROM sessions
		WHERE id = $1 AND expires_at > now()
	`

	var userID int64
	err := s.Pool.QueryRow(ctx, query, token).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, pgx.ErrNoRows
		}
		return 0, err
	}

	return userID, nil
}

// DeleteSession removes a session row for the supplied token.
func (s *Store) DeleteSession(ctx context.Context, token string) error {
	_, err := s.Pool.Exec(ctx, "DELETE FROM sessions WHERE id = $1", token)
	return err
}
