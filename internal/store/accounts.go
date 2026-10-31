package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// ErrResetInvalid is returned when a reset token is unknown, expired, or already used.
var ErrResetInvalid = errors.New("reset link is invalid or has expired")

// UpdateUserName changes the display name and returns nothing on success.
func (s *Store) UpdateUserName(ctx context.Context, userID int64, name string) error {
	_, err := s.Pool.Exec(ctx, `UPDATE users SET name = $2 WHERE id = $1`, userID, name)
	return err
}

// UpdatePasswordKeepingSession replaces the password hash and signs out every other session.
func (s *Store) UpdatePasswordKeepingSession(ctx context.Context, userID int64, hash, keepToken string) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `UPDATE users SET password_hash = $2 WHERE id = $1`, userID, hash); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1 AND id <> $2`, userID, keepToken); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// CreatePasswordReset stores the hash of a reset token. Older tokens for the user are removed.
func (s *Store) CreatePasswordReset(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM password_resets WHERE user_id = $1`, userID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO password_resets (token_hash, user_id, expires_at)
		VALUES ($1, $2, $3)
	`, tokenHash, userID, expiresAt); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// ConsumePasswordReset validates the token, sets the new password, and signs the user out everywhere.
// The token row is deleted in the same transaction so a link works exactly once.
func (s *Store) ConsumePasswordReset(ctx context.Context, tokenHash, newPasswordHash string) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var userID int64
	err = tx.QueryRow(ctx, `
		DELETE FROM password_resets
		WHERE token_hash = $1 AND expires_at > now()
		RETURNING user_id
	`, tokenHash).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrResetInvalid
		}
		return err
	}

	if _, err := tx.Exec(ctx, `UPDATE users SET password_hash = $2 WHERE id = $1`, userID, newPasswordHash); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// CreateContactMessage saves a contact form submission so it is never lost if email fails.
func (s *Store) CreateContactMessage(ctx context.Context, name, email, topic, message string) error {
	_, err := s.Pool.Exec(ctx, `
		INSERT INTO contact_messages (name, email, topic, message)
		VALUES ($1, $2, $3, $4)
	`, name, email, topic, message)
	return err
}

// DeleteExpired removes expired sessions and reset tokens.
func (s *Store) DeleteExpired(ctx context.Context) error {
	if _, err := s.Pool.Exec(ctx, `DELETE FROM sessions WHERE expires_at < now()`); err != nil {
		return err
	}
	_, err := s.Pool.Exec(ctx, `DELETE FROM password_resets WHERE expires_at < now()`)
	return err
}
