CREATE TABLE days (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL UNIQUE
);

CREATE TABLE habits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    start_date TEXT NOT NULL,
    archived_at TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE habit_logs (
    day_id INTEGER NOT NULL,
    habit_id INTEGER NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('completed', 'skipped')),
    completed_at TEXT,
    PRIMARY KEY (day_id, habit_id),
    FOREIGN KEY (day_id) REFERENCES days(id) ON DELETE CASCADE,
    FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
);

CREATE INDEX idx_habit_logs_habit_id ON habit_logs(habit_id);
CREATE INDEX idx_habits_active_dates ON habits(start_date, archived_at);
