package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

type DB struct{ SQL *sql.DB }

func Open(ctx context.Context, path string) (*DB, error) {
	if path != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, fmt.Errorf("mkdir db: %w", err)
		}
	}
	conn, err := sql.Open("sqlite3", path+"?_foreign_keys=on&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if path == ":memory:" {
		conn.SetMaxOpenConns(1)
		conn.SetMaxIdleConns(1)
	} else {
		conn.SetMaxOpenConns(8)
		conn.SetMaxIdleConns(8)
	}
	conn.SetConnMaxLifetime(time.Hour)
	db := &DB{SQL: conn}
	if err := db.Migrate(ctx); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return db, nil
}
func (db *DB) Close() error { return db.SQL.Close() }
func (db *DB) Migrate(ctx context.Context) error {
	if _, err := db.SQL.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS schema_migrations(version INTEGER PRIMARY KEY, applied_at TEXT NOT NULL)"); err != nil {
		return fmt.Errorf("migration table: %w", err)
	}
	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		content, readErr := migrationFS.ReadFile("migrations/" + entry.Name())
		if readErr != nil {
			return readErr
		}
		var version int
		if _, scanErr := fmt.Sscanf(entry.Name(), "%d_", &version); scanErr != nil {
			return scanErr
		}
		var exists int
		if err = db.SQL.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version).Scan(&exists); err != nil {
			return err
		}
		if exists > 0 {
			continue
		}
		tx, txErr := db.SQL.BeginTx(ctx, nil)
		if txErr != nil {
			return txErr
		}
		for _, statement := range strings.Split(string(content), ";") {
			statement = strings.TrimSpace(statement)
			if statement != "" {
				if _, txErr = tx.ExecContext(ctx, statement); txErr != nil {
					_ = tx.Rollback()
					return fmt.Errorf("migration %d: %w", version, txErr)
				}
			}
		}
		if _, txErr = tx.ExecContext(ctx, "INSERT INTO schema_migrations(version, applied_at) VALUES(?, ?)", version, time.Now().UTC().Format(time.RFC3339Nano)); txErr != nil {
			_ = tx.Rollback()
			return txErr
		}
		if txErr = tx.Commit(); txErr != nil {
			return txErr
		}
	}
	return nil
}
