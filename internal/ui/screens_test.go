package ui

import (
	"context"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"

	kaizenapp "kaizen-calendar/internal/app"
	"kaizen-calendar/internal/storage"
)

func TestScreensConstructWithCompletedItems(t *testing.T) {
	fyneApp := test.NewApp()
	defer fyneApp.Quit()

	db, err := storage.Open(filepath.Join(t.TempDir(), "screens.db"))
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	habitService := kaizenapp.NewHabitService(storage.NewHabitRepository(db))
	if err := habitService.Create(ctx, "Тестовая привычка", "Описание"); err != nil {
		t.Fatalf("create habit: %v", err)
	}
	today := startOfToday()
	habits, err := habitService.ListForDate(ctx, today)
	if err != nil {
		t.Fatalf("list habits: %v", err)
	}
	if err := habitService.SetStatus(ctx, today, habits[0].ID, storage.HabitStatusCompleted); err != nil {
		t.Fatalf("complete habit: %v", err)
	}

	todoService := kaizenapp.NewTodoService(storage.NewTodoRepository(db))
	if err := todoService.Create(ctx, "Тестовая задача"); err != nil {
		t.Fatalf("create todo: %v", err)
	}
	todos, err := todoService.List(ctx)
	if err != nil {
		t.Fatalf("list todos: %v", err)
	}
	if err := todoService.SetCompleted(ctx, todos[0].ID, true); err != nil {
		t.Fatalf("complete todo: %v", err)
	}

	window := test.NewWindow(container.NewWithoutLayout())
	defer window.Close()
	NewHabitScreen(habitService, window)
	NewStatisticsScreen(habitService)
	NewPomodoroScreen(window, true)
	NewTodoScreen(todoService, window)
}
