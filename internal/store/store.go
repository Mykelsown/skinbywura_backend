// Package store is the persistence layer for application data access.
// It is currently a placeholder for database and repository logic.
package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store holds the database connection pool used by data access methods.
type Store struct {
	Pool *pgxpool.Pool // Pool manages connections to PostgreSQL.
}

// Connection creates and validates a PostgreSQL connection pool.
func Connection(connectionString string) (*Store, error) {
	// Create the pool from the configured PostgreSQL connection string.
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	// Retry the first ping: on a fresh deploy the database may still be starting.
	var pingErr error
	for attempt := 1; attempt <= 10; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		pingErr = pool.Ping(ctx)
		cancel()
		if pingErr == nil {
			return &Store{Pool: pool}, nil
		}
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	pool.Close()
	return nil, pingErr
}
