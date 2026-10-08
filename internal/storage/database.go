package storage

import (
	"database/sql"
	"fmt"

	"kaizen-calendar/migrations"
	_ "modernc.org/sqlite"
)

// Open opens a local SQLite database and applies pending migrations.
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	// SQLite connections keep connection-local PRAGMA settings. A single
	// connection keeps foreign keys enabled consistently for this local app.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if _, err = db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}
	if _, err = db.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set busy timeout: %w", err)
	}
	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	if err = migrations.Apply(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply migrations: %w", err)
	}

	return db, nil
}
