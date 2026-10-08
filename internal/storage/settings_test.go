package storage

import (
	"context"
	"path/filepath"
	"testing"
)

func TestSettingsRepository(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "settings.db"))
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	defer db.Close()

	repository := NewSettingsRepository(db)
	ctx := context.Background()
	value, err := repository.Get(ctx, "missing", "fallback")
	if err != nil || value != "fallback" {
		t.Fatalf("unexpected fallback: value=%q err=%v", value, err)
	}
	if err := repository.Set(ctx, "sound", "false"); err != nil {
		t.Fatalf("set setting: %v", err)
	}
	value, err = repository.Get(ctx, "sound", "true")
	if err != nil || value != "false" {
		t.Fatalf("unexpected saved value: value=%q err=%v", value, err)
	}
}
