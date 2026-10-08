package app

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"kaizen-calendar/internal/storage"
)

func TestStatisticsCanUseHistoricalEndDate(t *testing.T) {
	db, err := storage.Open(filepath.Join(t.TempDir(), "statistics.db"))
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	defer db.Close()

	service := NewHabitService(storage.NewHabitRepository(db))
	ctx := context.Background()
	if err := service.Create(ctx, "Историческая привычка", ""); err != nil {
		t.Fatalf("create habit: %v", err)
	}

	endDate := time.Now().AddDate(0, 0, -7)
	statistics, err := service.StatisticsAt(ctx, PeriodWeek, endDate)
	if err != nil {
		t.Fatalf("get historical statistics: %v", err)
	}
	if !statistics.EndDate.Equal(time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())) {
		t.Fatalf("unexpected historical end date: %s", statistics.EndDate)
	}
}

func TestStatisticsAverageIgnoresDaysWithoutHabits(t *testing.T) {
	db, err := storage.Open(filepath.Join(t.TempDir(), "statistics.db"))
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	defer db.Close()

	service := NewHabitService(storage.NewHabitRepository(db))
	ctx := context.Background()

	// Anchor the test to a fixed mid-month date so the habit's start date
	// always falls inside the monthly statistics period regardless of when
	// the test runs.
	today := startOfDate(time.Date(2026, time.March, 15, 0, 0, 0, 0, time.Local))
	if _, err := service.repository.Create(ctx, "Новая привычка", "", today.AddDate(0, 0, -3).Format("2006-01-02")); err != nil {
		t.Fatalf("create habit: %v", err)
	}
	habitList, err := service.ListAll(ctx)
	if err != nil {
		t.Fatalf("list habits: %v", err)
	}
	for offset := -2; offset <= 0; offset++ {
		if err := service.SetStatus(ctx, today.AddDate(0, 0, offset), habitList[len(habitList)-1].ID, storage.HabitStatusCompleted); err != nil {
			t.Fatalf("complete habit on day %d: %v", offset, err)
		}
	}

	statistics, err := service.StatisticsAt(ctx, PeriodMonth, today)
	if err != nil {
		t.Fatalf("get monthly statistics: %v", err)
	}
	if statistics.ElapsedDays != today.Day() {
		t.Fatalf("unexpected elapsed days: %d", statistics.ElapsedDays)
	}
	if statistics.CompletedDays != 3 {
		t.Fatalf("unexpected completed days: %d", statistics.CompletedDays)
	}
	if statistics.AveragePercent != 75 {
		t.Fatalf("average should only consider days with habits, got %d", statistics.AveragePercent)
	}
	habitStats := statistics.Habits[len(statistics.Habits)-1]
	if habitStats.ApplicableDays != 4 || habitStats.CompletedDays != 3 || habitStats.Percent != 75 {
		t.Fatalf("unexpected habit totals: %#v", habitStats)
	}
}
