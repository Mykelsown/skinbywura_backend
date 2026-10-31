// Package migrate applies the embedded *.up.sql files at startup.
//
// It is intentionally compatible with golang-migrate: it uses the same
// schema_migrations(version, dirty) table, so databases already migrated with
// the `migrate` CLI during development keep working, and the CLI can still be used.
package migrate

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// advisoryLockKey stops two instances from migrating at the same moment.
const advisoryLockKey int64 = 7_251_001

// Migration is one numbered up file.
type Migration struct {
	Version int64
	Name    string
	SQL     string
}

// Load reads and orders the up migrations from dir inside fsys.
func Load(fsys fs.FS, dir string) ([]Migration, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}

	var out []Migration
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		v, err := ParseVersion(name)
		if err != nil {
			return nil, err
		}
		b, err := fs.ReadFile(fsys, path.Join(dir, name))
		if err != nil {
			return nil, err
		}
		out = append(out, Migration{Version: v, Name: name, SQL: string(b)})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Version < out[j].Version })
	for i := 1; i < len(out); i++ {
		if out[i].Version == out[i-1].Version {
			return nil, fmt.Errorf("duplicate migration version %d", out[i].Version)
		}
	}
	return out, nil
}

// ParseVersion extracts the numeric prefix of a name like 000007_create_x.up.sql.
func ParseVersion(name string) (int64, error) {
	prefix, _, ok := strings.Cut(name, "_")
	if !ok {
		return 0, fmt.Errorf("migration %q has no version prefix", name)
	}
	v, err := strconv.ParseInt(prefix, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("migration %q has a bad version prefix: %w", name, err)
	}
	return v, nil
}

// Pending returns the migrations newer than the current version.
func Pending(all []Migration, current int64) []Migration {
	var out []Migration
	for _, m := range all {
		if m.Version > current {
			out = append(out, m)
		}
	}
	return out
}

// Up applies every pending migration, each inside its own transaction.
func Up(ctx context.Context, pool *pgxpool.Pool, fsys fs.FS, dir string) error {
	all, err := Load(fsys, dir)
	if err != nil {
		return fmt.Errorf("load migrations: %w", err)
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, "SELECT pg_advisory_lock($1)", advisoryLockKey); err != nil {
		return err
	}
	defer conn.Exec(context.Background(), "SELECT pg_advisory_unlock($1)", advisoryLockKey)

	if _, err := conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version bigint NOT NULL PRIMARY KEY, dirty boolean NOT NULL)`); err != nil {
		return err
	}

	var current int64
	var dirty bool
	err = conn.QueryRow(ctx, `SELECT version, dirty FROM schema_migrations LIMIT 1`).Scan(&current, &dirty)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if dirty {
		return fmt.Errorf("database is at dirty version %d: fix it manually, then clear the dirty flag", current)
	}

	pending := Pending(all, current)
	if len(pending) == 0 {
		slog.Info("migrations up to date", "version", current)
		return nil
	}

	for _, m := range pending {
		tx, err := conn.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, m.SQL); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("apply %s: %w", m.Name, err)
		}
		if _, err := tx.Exec(ctx, `DELETE FROM schema_migrations`); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version, dirty) VALUES ($1, false)`, m.Version); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
		slog.Info("applied migration", "name", m.Name)
	}
	return nil
}
