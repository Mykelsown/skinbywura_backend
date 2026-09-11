package store

import (
	"context"
	"errors"

	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
	"github.com/jackc/pgx/v5"
)

// CreateUser inserts a new user row and returns the created record.
func (s *Store) CreateUser(ctx context.Context, name, email, passwordHash string) (types.User, error) {
	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, password_hash, created_at
	`

	var user types.User
	err := s.Pool.QueryRow(ctx, query, name, email, passwordHash).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}

// GetUserByEmail returns the matching user record, including the password hash.
func (s *Store) GetUserByEmail(ctx context.Context, email string) (types.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`

	var user types.User
	err := s.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.User{}, pgx.ErrNoRows
		}
		return types.User{}, err
	}

	return user, nil
}

// GetUserByID returns a user record by its identifier.
func (s *Store) GetUserByID(ctx context.Context, id int64) (types.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`

	var user types.User
	err := s.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.User{}, pgx.ErrNoRows
		}
		return types.User{}, err
	}

	return user, nil
}
