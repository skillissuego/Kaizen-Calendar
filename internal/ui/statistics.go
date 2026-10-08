package ui

import (
	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	kaizenapp "kaizen-calendar/internal/app"
	"kaizen-calendar/internal/storage"
)

type StatisticsScreen struct {
	service      *kaizenapp.HabitService
	period       kaizenapp.StatisticsPeriod
	periodSelect *widget.Select
	endDate      *widget.DateEntry
	selectedEnd  time.Time
	summary      *widget.Label
	rows         *fyne.Container
	root         *fyne.Container
}

func NewStatisticsScreen(service *kaizenapp.HabitService) fyne.CanvasObject {
	screen := &StatisticsScreen{
		service:      service,
		period:       kaizenapp.PeriodMonth,
		periodSelect: widget.NewSelect([]string{"Неделя", "Месяц", "Год"}, nil),
		endDate:      widget.NewDateEntry(),
		selectedEnd:  startOfToday(),
		summary:      widget.NewLabel(""),
		rows:         container.NewVBox(),
	}
	screen.periodSelect.OnChanged = func(selected string) {
		switch selected {
		case "Неделя":
			screen.period = kaizenapp.PeriodWeek
		case "Год":
			screen.period = kaizenapp.PeriodYear
		default:
			screen.period = kaizenapp.PeriodMonth
		}
		screen.refresh()
	}
	screen.endDate.OnChanged = func(selected *time.Time) {
		if selected == nil {
			screen.selectedEnd = startOfToday()
		} else {
			screen.selectedEnd = *selected
		}
		if screen.selectedEnd.After(startOfToday()) {
			screen.selectedEnd = startOfToday()
		}
		screen.refresh()
	}
	screen.periodSelect.SetSelected("Месяц")
	date := screen.selectedEnd
	screen.endDate.SetDate(&date)

	entryHeight := screen.endDate.MinSize().Height
	dateWrap := container.New(layout.NewGridWrapLayout(fyne.NewSize(180, entryHeight)), screen.endDate)
	periodRow := container.NewHBox(widget.NewLabel("Период:"), screen.periodSelect)
	dateRow := container.NewHBox(widget.NewLabel("До даты:"), dateWrap)
	header := container.NewVBox(periodRow, dateRow)
	content := container.NewBorder(
		header,
		nil,
		nil,
		nil,
		container.NewVScroll(screen.rows),
	)
	screen.root = container.NewBorder(screen.summary, nil, nil, nil, content)
	screen.refresh()
	return screen.root
}

func (s *StatisticsScreen) Refresh() {
	s.refresh()
}

func (s *StatisticsScreen) refresh() {
	statistics, err := s.service.StatisticsAt(context.Background(), s.period, s.selectedEnd)
	if err != nil {
		s.summary.SetText("Ошибка загрузки статистики: " + err.Error())
		s.rows.Objects = nil
		s.rows.Refresh()
		return
	}

	s.summary.SetText(fmt.Sprintf(
		"Статистика за период: %s - %s\n"+
			"Полностью завершённых дней: %d\n"+
			"Среднее выполнение дня: %d%%",
		statistics.StartDate.Format("02.01.2006"),
		statistics.EndDate.Format("02.01.2006"),
		statistics.CompletedDays,
		statistics.AveragePercent,
	))

	s.rows.Objects = nil
	if len(statistics.Habits) == 0 {
		s.rows.Add(widget.NewLabel("Добавьте первую привычку, чтобы появилась статистика."))
		s.rows.Refresh()
		return
	}
	s.rows.Add(widget.NewSeparator())
	for _, habit := range statistics.Habits {
		s.rows.Add(statisticsRow(habit))
	}
	s.rows.Refresh()
}

func statisticsRow(statistics storage.HabitStatistics) fyne.CanvasObject {
	name := widget.NewLabel(statistics.Habit.Name)
	name.TextStyle = fyne.TextStyle{Bold: true}

	status := widget.NewLabel("Активна")
	if statistics.Habit.ArchivedAt.Valid {
		status.SetText("Архивирована")
		status.Importance = widget.LowImportance
	}
	nameRow := container.NewHBox(name, status)

	details := widget.NewLabel(fmt.Sprintf(
		"Выполнено: %d из %d дней (%d%%)   Серия: %d   Максимум: %d",
		statistics.CompletedDays,
		statistics.ApplicableDays,
		statistics.Percent,
		statistics.CurrentStreak,
		statistics.MaxStreak,
	))

	items := []fyne.CanvasObject{nameRow, details, widget.NewSeparator()}
	return container.NewVBox(items...)
}
