package storage

import (
	"context"
	"fmt"
	"time"
)

type HabitStatistics struct {
	Habit          Habit
	ApplicableDays int
	CompletedDays  int
	Percent        int
	CurrentStreak  int
	MaxStreak      int
}

func (r *HabitRepository) GetHabitStatistics(ctx context.Context, habitID int64, startDate, endDate string) (HabitStatistics, error) {
	var result HabitStatistics
	if err := r.db.QueryRowContext(ctx, `
		SELECT id, name, description, start_date, archived_at, sort_order, created_at
		FROM habits WHERE id = ?
	`, habitID).Scan(
		&result.Habit.ID,
		&result.Habit.Name,
		&result.Habit.Description,
		&result.Habit.StartDate,
		&result.Habit.ArchivedAt,
		&result.Habit.SortOrder,
		&result.Habit.CreatedAt,
	); err != nil {
		return HabitStatistics{}, fmt.Errorf("load habit statistics: %w", err)
	}

	periodStart, err := parseDate(startDate)
	if err != nil {
		return HabitStatistics{}, fmt.Errorf("parse statistics start date: %w", err)
	}
	periodEnd, err := parseDate(endDate)
	if err != nil {
		return HabitStatistics{}, fmt.Errorf("parse statistics end date: %w", err)
	}
	effectiveStart, err := parseDate(result.Habit.StartDate)
	if err != nil {
		return HabitStatistics{}, fmt.Errorf("parse habit start date: %w", err)
	}
	if effectiveStart.Before(periodStart) {
		effectiveStart = periodStart
	}
	effectiveEnd := periodEnd
	if result.Habit.ArchivedAt.Valid {
		archiveDate, parseErr := parseDate(result.Habit.ArchivedAt.String)
		if parseErr != nil {
			return HabitStatistics{}, fmt.Errorf("parse habit archive date: %w", parseErr)
		}
		if archiveDate.Before(effectiveEnd) {
			effectiveEnd = archiveDate
		}
	}

	if !effectiveStart.After(effectiveEnd) {
		result.ApplicableDays = int(effectiveEnd.Sub(effectiveStart).Hours()/24) + 1
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT d.date
		FROM habit_logs l
		JOIN days d ON d.id = l.day_id
		WHERE l.habit_id = ?
		  AND l.status = 'completed'
		  AND d.date BETWEEN ? AND ?
		ORDER BY d.date
	`, habitID, startDate, endDate)
	if err != nil {
		return HabitStatistics{}, fmt.Errorf("load completed habit dates: %w", err)
	}
	defer rows.Close()

	var completedDates []time.Time
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			return HabitStatistics{}, fmt.Errorf("scan completed habit date: %w", err)
		}
		parsed, err := parseDate(date)
		if err != nil {
			return HabitStatistics{}, fmt.Errorf("parse completed habit date: %w", err)
		}
		if !parsed.Before(effectiveStart) && !parsed.After(effectiveEnd) {
			completedDates = append(completedDates, parsed)
		}
	}
	if err := rows.Err(); err != nil {
		return HabitStatistics{}, fmt.Errorf("iterate completed habit dates: %w", err)
	}

	result.CompletedDays = len(completedDates)
	if result.ApplicableDays > 0 {
		result.Percent = result.CompletedDays * 100 / result.ApplicableDays
	}
	result.CurrentStreak, result.MaxStreak = calculateStreaks(completedDates, effectiveEnd)
	return result, nil
}

func calculateStreaks(dates []time.Time, periodEnd time.Time) (int, int) {
	if len(dates) == 0 {
		return 0, 0
	}

	maxStreak := 1
	currentRun := 1
	for index := 1; index < len(dates); index++ {
		if dates[index].Sub(dates[index-1]) == 24*time.Hour {
			currentRun++
		} else {
			if currentRun > maxStreak {
				maxStreak = currentRun
			}
			currentRun = 1
		}
	}
	if currentRun > maxStreak {
		maxStreak = currentRun
	}

	currentStreak := 0
	if dates[len(dates)-1].Equal(periodEnd) {
		currentStreak = 1
		for index := len(dates) - 1; index > 0; index-- {
			if dates[index].Sub(dates[index-1]) != 24*time.Hour {
				break
			}
			currentStreak++
		}
	}
	return currentStreak, maxStreak
}

func parseDate(value string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", value, time.Local)
}
