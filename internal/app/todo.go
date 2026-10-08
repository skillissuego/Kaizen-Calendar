package app

import (
	"context"
	"fmt"
	"strings"

	"kaizen-calendar/internal/storage"
)

type TodoService struct {
	repository *storage.TodoRepository
}

func NewTodoService(repository *storage.TodoRepository) *TodoService {
	return &TodoService{repository: repository}
}

func (s *TodoService) Create(ctx context.Context, title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return fmt.Errorf("название задачи не может быть пустым")
	}
	return s.repository.Create(ctx, title)
}

func (s *TodoService) List(ctx context.Context) ([]storage.TodoItem, error) {
	return s.repository.List(ctx)
}

func (s *TodoService) SetCompleted(ctx context.Context, id int64, completed bool) error {
	return s.repository.SetCompleted(ctx, id, completed)
}

func (s *TodoService) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
