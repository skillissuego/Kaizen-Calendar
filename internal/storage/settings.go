package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Get(ctx context.Context, key, fallback string) (string, error) {
	var value string
	err := r.db.QueryRowContext(ctx, "SELECT value FROM app_settings WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return fallback, nil
	}
	if err != nil {
		return "", fmt.Errorf("get setting %q: %w", key, err)
	}
	return value, nil
}

func (r *SettingsRepository) Set(ctx context.Context, key, value string) error {
	if _, err := r.db.ExecContext(ctx, `
		INSERT INTO app_settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, key, value); err != nil {
		return fmt.Errorf("set setting %q: %w", key, err)
	}
	return nil
}
