package storage

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestHabitLifecycleAndDateRules(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	repository := NewHabitRepository(db)
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	habitID, err := repository.Create(ctx, "Чтение", "", today)
	if err != nil {
		t.Fatalf("create habit: %v", err)
	}

	yesterdayHabits, err := repository.ListForDate(ctx, yesterday)
	if err != nil {
		t.Fatalf("list yesterday habits: %v", err)
	}
	if len(yesterdayHabits) != 0 {
		t.Fatalf("new habit appeared before start date: %#v", yesterdayHabits)
	}

	dayID, err := repository.EnsureDay(ctx, today)
	if err != nil {
		t.Fatalf("ensure today: %v", err)
	}
	if err := repository.SetStatus(ctx, dayID, habitID, HabitStatusCompleted); err != nil {
		t.Fatalf("set completed status: %v", err)
	}

	todayHabits, err := repository.ListForDate(ctx, today)
	if err != nil {
		t.Fatalf("list today habits: %v", err)
	}
	if len(todayHabits) != 1 || todayHabits[0].Status != HabitStatusCompleted {
		t.Fatalf("unexpected today habit status: %#v", todayHabits)
	}
	progress, err := repository.GetDayProgress(ctx, today, 0)
	if err != nil {
		t.Fatalf("get today's progress: %v", err)
	}
	if progress.Completed != 1 || progress.Total != 1 || progress.Percent != 100 {
		t.Fatalf("unexpected progress: %#v", progress)
	}

	if err := repository.SetStatus(ctx, dayID, habitID, HabitStatusSkipped); err != nil {
		t.Fatalf("set skipped status: %v", err)
	}
	progress, err = repository.GetDayProgress(ctx, today, 0)
	if err != nil {
		t.Fatalf("get skipped progress: %v", err)
	}
	if progress.Completed != 0 || progress.Total != 0 || progress.Percent != 0 {
		t.Fatalf("unexpected skipped progress: %#v", progress)
	}

	if err := repository.Archive(ctx, habitID, today); err != nil {
		t.Fatalf("archive habit: %v", err)
	}
	tomorrowHabits, err := repository.ListForDate(ctx, tomorrow)
	if err != nil {
		t.Fatalf("list tomorrow habits: %v", err)
	}
	if len(tomorrowHabits) != 0 {
		t.Fatalf("archived habit appeared after archive date: %#v", tomorrowHabits)
	}
}

func TestHabitStatisticsCalculateStreaks(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "statistics.db"))
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	repository := NewHabitRepository(db)
	today := time.Now()
	start := today.AddDate(0, 0, -4).Format("2006-01-02")
	habitID, err := repository.Create(ctx, "Статистика", "", start)
	if err != nil {
		t.Fatalf("create habit: %v", err)
	}

	for offset := -3; offset <= 0; offset++ {
		date := today.AddDate(0, 0, offset).Format("2006-01-02")
		dayID, err := repository.EnsureDay(ctx, date)
		if err != nil {
			t.Fatalf("ensure day %s: %v", date, err)
		}
		if err := repository.SetStatus(ctx, dayID, habitID, HabitStatusCompleted); err != nil {
			t.Fatalf("complete habit for %s: %v", date, err)
		}
	}

	statistics, err := repository.GetHabitStatistics(
		ctx,
		habitID,
		today.AddDate(0, 0, -4).Format("2006-01-02"),
		today.Format("2006-01-02"),
	)
	if err != nil {
		t.Fatalf("get habit statistics: %v", err)
	}
	if statistics.CompletedDays != 4 || statistics.ApplicableDays != 5 || statistics.Percent != 80 {
		t.Fatalf("unexpected totals: %#v", statistics)
	}
	if statistics.CurrentStreak != 4 || statistics.MaxStreak != 4 {
		t.Fatalf("unexpected streaks: %#v", statistics)
	}
}
