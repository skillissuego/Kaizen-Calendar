package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type TodoItem struct {
	ID          int64
	Title       string
	Completed   bool
	CreatedAt   string
	CompletedAt sql.NullString
}

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(ctx context.Context, title string) error {
	if _, err := r.db.ExecContext(ctx, "INSERT INTO todo_items (title) VALUES (?)", title); err != nil {
		return fmt.Errorf("create todo item: %w", err)
	}
	return nil
}

func (r *TodoRepository) List(ctx context.Context) ([]TodoItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, completed, created_at, completed_at
		FROM todo_items
		ORDER BY completed, created_at DESC, id DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list todo items: %w", err)
	}
	defer rows.Close()

	var items []TodoItem
	for rows.Next() {
		var item TodoItem
		var completed int
		if err := rows.Scan(&item.ID, &item.Title, &completed, &item.CreatedAt, &item.CompletedAt); err != nil {
			return nil, fmt.Errorf("scan todo item: %w", err)
		}
		item.Completed = completed == 1
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate todo items: %w", err)
	}
	return items, nil
}

func (r *TodoRepository) SetCompleted(ctx context.Context, id int64, completed bool) error {
	completedValue := 0
	if completed {
		completedValue = 1
	}
	if _, err := r.db.ExecContext(ctx, `
		UPDATE todo_items
		SET completed = ?, completed_at = CASE WHEN ? = 1 THEN CURRENT_TIMESTAMP ELSE NULL END
		WHERE id = ?
	`, completedValue, completedValue, id); err != nil {
		return fmt.Errorf("set todo item status: %w", err)
	}
	return nil
}

func (r *TodoRepository) Delete(ctx context.Context, id int64) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM todo_items WHERE id = ?", id); err != nil {
		return fmt.Errorf("delete todo item: %w", err)
	}
	return nil
}
