package storage

import (
	"context"
	"path/filepath"
	"testing"
)

func TestTodoLifecycle(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "todo.db"))
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	repository := NewTodoRepository(db)
	if err := repository.Create(ctx, "Проверить Todo"); err != nil {
		t.Fatalf("create todo item: %v", err)
	}

	items, err := repository.List(ctx)
	if err != nil {
		t.Fatalf("list todo items: %v", err)
	}
	if len(items) != 1 || items[0].Completed {
		t.Fatalf("unexpected todo items: %#v", items)
	}

	if err := repository.SetCompleted(ctx, items[0].ID, true); err != nil {
		t.Fatalf("complete todo item: %v", err)
	}
	items, err = repository.List(ctx)
	if err != nil {
		t.Fatalf("list completed todo items: %v", err)
	}
	if len(items) != 1 || !items[0].Completed {
		t.Fatalf("todo item was not completed: %#v", items)
	}

	if err := repository.Delete(ctx, items[0].ID); err != nil {
		t.Fatalf("delete todo item: %v", err)
	}
	items, err = repository.List(ctx)
	if err != nil {
		t.Fatalf("list deleted todo items: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("todo item was not deleted: %#v", items)
	}
}
