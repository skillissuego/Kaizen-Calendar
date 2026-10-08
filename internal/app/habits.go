package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"kaizen-calendar/internal/storage"
)

type HabitService struct {
	repository *storage.HabitRepository
}

func NewHabitService(repository *storage.HabitRepository) *HabitService {
	return &HabitService{repository: repository}
}

func (s *HabitService) Create(ctx context.Context, name, description string) error {
	name = normalizeName(name)
	if name == "" {
		return fmt.Errorf("название привычки не может быть пустым")
	}
	_, err := s.repository.Create(ctx, name, description, today())
	return err
}

func (s *HabitService) Update(ctx context.Context, id int64, name, description string) error {
	name = normalizeName(name)
	if name == "" {
		return fmt.Errorf("название привычки не может быть пустым")
	}
	return s.repository.Update(ctx, id, name, description)
}

func (s *HabitService) Archive(ctx context.Context, id int64) error {
	return s.repository.Archive(ctx, id, today())
}

func (s *HabitService) DeleteArchived(ctx context.Context, id int64) error {
	return s.repository.DeleteArchived(ctx, id)
}

func (s *HabitService) ListAll(ctx context.Context) ([]storage.Habit, error) {
	return s.repository.ListAll(ctx)
}

func (s *HabitService) ListForDate(ctx context.Context, date time.Time) ([]storage.DailyHabit, error) {
	return s.repository.ListForDate(ctx, formatDate(date))
}

func (s *HabitService) GetDayProgress(ctx context.Context, date time.Time, habitID int64) (storage.DayProgress, error) {
	return s.repository.GetDayProgress(ctx, formatDate(date), habitID)
}

func (s *HabitService) SetStatus(ctx context.Context, date time.Time, habitID int64, status string) error {
	if date.After(startOfToday()) {
		return fmt.Errorf("будущие даты нельзя изменять")
	}
	dayID, err := s.repository.EnsureDay(ctx, formatDate(date))
	if err != nil {
		return err
	}
	return s.repository.SetStatus(ctx, dayID, habitID, status)
}

func (s *HabitService) ClearStatus(ctx context.Context, date time.Time, habitID int64) error {
	if date.After(startOfToday()) {
		return fmt.Errorf("будущие даты нельзя изменять")
	}
	dayID, err := s.repository.EnsureDay(ctx, formatDate(date))
	if err != nil {
		return err
	}
	return s.repository.ClearStatus(ctx, dayID, habitID)
}

func (s *HabitService) Reorder(ctx context.Context, ids []int64) error {
	return s.repository.SetOrder(ctx, ids)
}

func formatDate(value time.Time) string {
	return value.Format("2006-01-02")
}

func today() string {
	return formatDate(time.Now())
}

func startOfToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func normalizeName(value string) string {
	return strings.TrimSpace(value)
}
