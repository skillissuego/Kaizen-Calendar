package app

import (
	"context"
	"fmt"
	"time"

	"kaizen-calendar/internal/storage"
)

type StatisticsPeriod string

const (
	PeriodWeek  StatisticsPeriod = "week"
	PeriodMonth StatisticsPeriod = "month"
	PeriodYear  StatisticsPeriod = "year"
)

type PeriodStatistics struct {
	StartDate      time.Time
	EndDate        time.Time
	CompletedDays  int
	ElapsedDays    int
	AveragePercent int
	Habits         []storage.HabitStatistics
}

func (s *HabitService) Statistics(ctx context.Context, period StatisticsPeriod, now time.Time) (PeriodStatistics, error) {
	return s.StatisticsAt(ctx, period, now)
}

func (s *HabitService) StatisticsAt(ctx context.Context, period StatisticsPeriod, endDate time.Time) (PeriodStatistics, error) {
	endDate = startOfDate(endDate)
	today := startOfDate(time.Now())
	if endDate.After(today) {
		endDate = today
	}
	start, end, err := periodRange(period, endDate)
	if err != nil {
		return PeriodStatistics{}, err
	}

	allHabits, err := s.repository.ListAll(ctx)
	if err != nil {
		return PeriodStatistics{}, fmt.Errorf("list habits for statistics: %w", err)
	}
	result := PeriodStatistics{
		StartDate: start,
		EndDate:   end,
	}
	for _, habit := range allHabits {
		habitStats, err := s.repository.GetHabitStatistics(ctx, habit.ID, formatDate(start), formatDate(end))
		if err != nil {
			return PeriodStatistics{}, err
		}
		result.Habits = append(result.Habits, habitStats)
	}

	percentTotal := 0
	daysWithHabits := 0
	elapsedDays := 0
	for date := start; !date.After(end); date = date.AddDate(0, 0, 1) {
		elapsedDays++
		progress, err := s.repository.GetDayProgress(ctx, formatDate(date), 0)
		if err != nil {
			return PeriodStatistics{}, fmt.Errorf("calculate daily summary: %w", err)
		}
		if progress.Total > 0 {
			percentTotal += progress.Percent
			daysWithHabits++
			if progress.Percent == 100 {
				result.CompletedDays++
			}
		}
	}
	result.ElapsedDays = elapsedDays
	if daysWithHabits > 0 {
		result.AveragePercent = percentTotal / daysWithHabits
	}
	return result, nil
}

func periodRange(period StatisticsPeriod, now time.Time) (time.Time, time.Time, error) {
	end := startOfDate(now)
	switch period {
	case PeriodWeek:
		mondayOffset := (int(end.Weekday()) + 6) % 7
		return end.AddDate(0, 0, -mondayOffset), end, nil
	case PeriodMonth:
		return time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, end.Location()), end, nil
	case PeriodYear:
		return time.Date(end.Year(), time.January, 1, 0, 0, 0, 0, end.Location()), end, nil
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("unknown statistics period %q", period)
	}
}

func startOfDate(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}
