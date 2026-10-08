package storage

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	HabitStatusCompleted = "completed"
	HabitStatusSkipped   = "skipped"
)

type Habit struct {
	ID          int64
	Name        string
	Description string
	StartDate   string
	ArchivedAt  sql.NullString
	SortOrder   int
	CreatedAt   string
}

type DailyHabit struct {
	Habit
	Status string
}

type DayProgress struct {
	Date      string
	Completed int
	Total     int
	Percent   int
}

type HabitRepository struct {
	db *sql.DB
}

func NewHabitRepository(db *sql.DB) *HabitRepository {
	return &HabitRepository{db: db}
}

func (r *HabitRepository) Create(ctx context.Context, name, description, startDate string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO habits (name, description, start_date, sort_order)
		VALUES (?, ?, ?, COALESCE((SELECT MAX(sort_order) + 1 FROM habits), 0))
	`, name, description, startDate)
	if err != nil {
		return 0, fmt.Errorf("create habit: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get created habit id: %w", err)
	}
	return id, nil
}

func (r *HabitRepository) Update(ctx context.Context, id int64, name, description string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE habits SET name = ?, description = ? WHERE id = ?
	`, name, description, id)
	if err != nil {
		return fmt.Errorf("update habit: %w", err)
	}

	if affected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("check updated habit: %w", err)
	} else if affected == 0 {
		return fmt.Errorf("habit %d not found", id)
	}
	return nil
}

func (r *HabitRepository) Archive(ctx context.Context, id int64, archiveDate string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE habits SET archived_at = ? WHERE id = ? AND archived_at IS NULL
	`, archiveDate, id)
	if err != nil {
		return fmt.Errorf("archive habit: %w", err)
	}
	if affected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("check archived habit: %w", err)
	} else if affected == 0 {
		return fmt.Errorf("active habit %d not found", id)
	}
	return nil
}

func (r *HabitRepository) DeleteArchived(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM habits WHERE id = ? AND archived_at IS NOT NULL
	`, id)
	if err != nil {
		return fmt.Errorf("delete archived habit: %w", err)
	}
	if affected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("check deleted habit: %w", err)
	} else if affected == 0 {
		return fmt.Errorf("archived habit %d not found", id)
	}
	return nil
}

func (r *HabitRepository) ListAll(ctx context.Context) ([]Habit, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, description, start_date, archived_at, sort_order, created_at
		FROM habits
		ORDER BY sort_order, name
	`)
	if err != nil {
		return nil, fmt.Errorf("list habits: %w", err)
	}
	defer rows.Close()

	var habits []Habit
	for rows.Next() {
		var habit Habit
		if err := rows.Scan(
			&habit.ID,
			&habit.Name,
			&habit.Description,
			&habit.StartDate,
			&habit.ArchivedAt,
			&habit.SortOrder,
			&habit.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan habit: %w", err)
		}
		habits = append(habits, habit)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate habits: %w", err)
	}
	return habits, nil
}

func (r *HabitRepository) ListForDate(ctx context.Context, date string) ([]DailyHabit, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT h.id, h.name, h.description, h.start_date, h.archived_at,
		       h.sort_order, h.created_at, COALESCE(l.status, '')
		FROM habits h
		LEFT JOIN days d ON d.date = ?
		LEFT JOIN habit_logs l ON l.habit_id = h.id AND l.day_id = d.id
		WHERE h.start_date <= ?
		  AND (h.archived_at IS NULL OR h.archived_at >= ?)
		ORDER BY h.sort_order, h.name
	`, date, date, date)
	if err != nil {
		return nil, fmt.Errorf("list daily habits: %w", err)
	}
	defer rows.Close()

	var habits []DailyHabit
	for rows.Next() {
		var habit DailyHabit
		if err := rows.Scan(
			&habit.ID,
			&habit.Name,
			&habit.Description,
			&habit.StartDate,
			&habit.ArchivedAt,
			&habit.SortOrder,
			&habit.CreatedAt,
			&habit.Status,
		); err != nil {
			return nil, fmt.Errorf("scan daily habit: %w", err)
		}
		habits = append(habits, habit)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate daily habits: %w", err)
	}
	return habits, nil
}

func (r *HabitRepository) GetDayProgress(ctx context.Context, date string, habitID int64) (DayProgress, error) {
	progress := DayProgress{Date: date}
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN l.status = 'completed' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN COALESCE(l.status, '') != 'skipped' THEN 1 ELSE 0 END), 0)
		FROM habits h
		LEFT JOIN days d ON d.date = ?
		LEFT JOIN habit_logs l ON l.habit_id = h.id AND l.day_id = d.id
		WHERE h.start_date <= ?
		  AND (h.archived_at IS NULL OR h.archived_at >= ?)
	`
	args := []any{date, date, date}
	if habitID > 0 {
		query += " AND h.id = ?"
		args = append(args, habitID)
	}

	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&progress.Completed, &progress.Total); err != nil {
		return DayProgress{}, fmt.Errorf("get progress for %s: %w", date, err)
	}
	if progress.Total > 0 {
		progress.Percent = progress.Completed * 100 / progress.Total
	}
	return progress, nil
}

func (r *HabitRepository) SetStatus(ctx context.Context, dayID, habitID int64, status string) error {
	if status != HabitStatusCompleted && status != HabitStatusSkipped {
		return fmt.Errorf("unsupported habit status %q", status)
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO habit_logs (day_id, habit_id, status, completed_at)
		VALUES (?, ?, ?, CASE WHEN ? = 'completed' THEN CURRENT_TIMESTAMP ELSE NULL END)
		ON CONFLICT(day_id, habit_id) DO UPDATE SET
			status = excluded.status,
			completed_at = excluded.completed_at
	`, dayID, habitID, status, status)
	if err != nil {
		return fmt.Errorf("set habit status: %w", err)
	}
	return nil
}

func (r *HabitRepository) ClearStatus(ctx context.Context, dayID, habitID int64) error {
	if _, err := r.db.ExecContext(ctx, `
		DELETE FROM habit_logs WHERE day_id = ? AND habit_id = ?
	`, dayID, habitID); err != nil {
		return fmt.Errorf("clear habit status: %w", err)
	}
	return nil
}

func (r *HabitRepository) EnsureDay(ctx context.Context, date string) (int64, error) {
	if _, err := r.db.ExecContext(ctx, `
		INSERT INTO days (date) VALUES (?) ON CONFLICT(date) DO NOTHING
	`, date); err != nil {
		return 0, fmt.Errorf("ensure day: %w", err)
	}

	var id int64
	if err := r.db.QueryRowContext(ctx, "SELECT id FROM days WHERE date = ?", date).Scan(&id); err != nil {
		return 0, fmt.Errorf("get day id: %w", err)
	}
	return id, nil
}

func (r *HabitRepository) SetOrder(ctx context.Context, ids []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin reorder habits: %w", err)
	}
	defer tx.Rollback()

	for order, id := range ids {
		if _, err := tx.ExecContext(ctx, "UPDATE habits SET sort_order = ? WHERE id = ?", order, id); err != nil {
			return fmt.Errorf("set habit order: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit reorder habits: %w", err)
	}
	return nil
}
