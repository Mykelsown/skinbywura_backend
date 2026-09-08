// Package store is the persistence layer for application data access.
// It is currently a placeholder for database and repository logic.
package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store holds the database connection pool used by data access methods.
type Store struct {
	Pool *pgxpool.Pool
}

// Connection creates and validates a PostgreSQL connection pool.
func Connection(connectionString string) (*Store, error) {
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return &Store{Pool: pool}, nil
}
