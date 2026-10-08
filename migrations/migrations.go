package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

//go:embed *.sql
var files embed.FS

// Apply applies each migration exactly once and records it in the database.
func Apply(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	names, err := fs.Glob(files, "*.sql")
	if err != nil {
		return fmt.Errorf("find migrations: %w", err)
	}
	sort.Strings(names)

	for _, name := range names {
		version, err := migrationVersion(name)
		if err != nil {
			return err
		}

		var applied bool
		err = db.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = ?)",
			version,
		).Scan(&applied)
		if err != nil {
			return fmt.Errorf("check migration %s: %w", name, err)
		}
		if applied {
			continue
		}

		sqlBytes, err := files.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin migration %s: %w", name, err)
		}

		if _, err = tx.Exec(string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("execute migration %s: %w", name, err)
		}
		if _, err = tx.Exec(
			"INSERT INTO schema_migrations (version, name) VALUES (?, ?)",
			version,
			name,
		); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", name, err)
		}
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", name, err)
		}
	}

	return nil
}

func migrationVersion(name string) (int, error) {
	base := filepath.Base(name)
	parts := strings.SplitN(base, "_", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid migration filename %q", name)
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil || version < 1 {
		return 0, fmt.Errorf("invalid migration version in %q", name)
	}
	return version, nil
}
