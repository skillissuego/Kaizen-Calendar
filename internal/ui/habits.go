package ui

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	kaizenapp "kaizen-calendar/internal/app"
	"kaizen-calendar/internal/storage"
)

type HabitScreen struct {
	service       *kaizenapp.HabitService
	window        fyne.Window
	selectedDate  time.Time
	dateLabel     *widget.Label
	dailyBox      *fyne.Container
	root          fyne.CanvasObject
	calendarRoot  *fyne.Container
	calendarSplit *container.Split
	nextDay       *widget.Button
	nextMonth     *widget.Button
	month         time.Time
	habitFilter   int64
	filterSelect  *widget.Select
	filterIDs     map[string]int64
	updating      bool
}

func NewHabitScreen(service *kaizenapp.HabitService, window fyne.Window) fyne.CanvasObject {
	screen := &HabitScreen{
		service:      service,
		window:       window,
		selectedDate: startOfToday(),
		dateLabel:    widget.NewLabel(""),
		dailyBox:     container.NewVBox(),
		calendarRoot: container.NewVBox(),
		month:        startOfMonth(startOfToday()),
	}

	previous := widget.NewButton("<", func() {
		screen.selectedDate = screen.selectedDate.AddDate(0, 0, -1)
		screen.refresh()
	})
	todayButton := widget.NewButton("Сегодня", func() {
		screen.selectedDate = startOfToday()
		screen.refresh()
	})
	next := widget.NewButton(">", func() {
		if !screen.selectedDate.Before(startOfToday()) {
			return
		}
		screen.selectedDate = screen.selectedDate.AddDate(0, 0, 1)
		screen.refresh()
	})
	screen.nextDay = next

	manage := widget.NewButton("Управление", screen.showHabitManager)
	add := widget.NewButton("Добавить", func() {
		screen.showHabitForm(nil)
	})

	// Stack the date label above the navigation row so the panel's minimum
	// width stays narrow and the calendar keeps most of the window.
	header := container.NewVBox(
		screen.dateLabel,
		container.NewHBox(previous, todayButton, next),
	)
	screen.dateLabel.Wrapping = fyne.TextWrapWord
	bottom := container.NewVBox(add, manage)
	left := container.NewBorder(header, bottom, nil, nil, container.NewVScroll(screen.dailyBox))
	screen.filterSelect = widget.NewSelect(nil, func(selected string) {
		screen.habitFilter = screen.filterIDs[selected]
		screen.refreshCalendar()
	})
	right := screen.calendarRoot
	split := container.NewHSplit(left, right)
	// The calendar is the primary visual surface, so it gets most of the width.
	split.SetOffset(0.18)
	screen.calendarSplit = split
	screen.root = split
	screen.refresh()
	return screen.root
}

func (s *HabitScreen) refresh() {
	s.dateLabel.SetText(fmt.Sprintf("Привычки на %s", russianDate(s.selectedDate)))
	if s.nextDay != nil {
		if s.selectedDate.Before(startOfToday()) {
			s.nextDay.Enable()
		} else {
			s.nextDay.Disable()
		}
	}
	s.refreshDaily()
	s.refreshCalendar()
}

func (s *HabitScreen) refreshDaily() {
	s.dailyBox.Objects = nil

	habits, err := s.service.ListForDate(context.Background(), s.selectedDate)
	if err != nil {
		s.dailyBox.Add(widget.NewLabel("Ошибка загрузки привычек: " + err.Error()))
		s.dailyBox.Refresh()
		return
	}
	if len(habits) == 0 {
		empty := widget.NewLabel("Пока нет привычек. Добавьте первую привычку.")
		empty.Wrapping = fyne.TextWrapWord
		s.dailyBox.Add(empty)
		s.dailyBox.Refresh()
		return
	}

	completed := 0
	considered := 0
	for _, habit := range habits {
		if habit.Status == storage.HabitStatusCompleted {
			completed++
		}
		if habit.Status != storage.HabitStatusSkipped {
			considered++
		}
		s.dailyBox.Add(s.dailyHabitRow(habit))
	}
	percent := 0
	if considered > 0 {
		percent = completed * 100 / considered
	}
	s.dailyBox.Add(widget.NewSeparator())
	summary := widget.NewLabel(fmt.Sprintf("Выполнено: %d из %d (%d%%)", completed, considered, percent))
	summary.Wrapping = fyne.TextWrapWord
	s.dailyBox.Add(summary)
	s.dailyBox.Refresh()
}

func (s *HabitScreen) refreshCalendar() {
	habits, err := s.service.ListAll(context.Background())
	if err != nil {
		s.calendarRoot.Objects = []fyne.CanvasObject{widget.NewLabel("Ошибка загрузки календаря: " + err.Error())}
		s.calendarRoot.Refresh()
		return
	}

	s.filterIDs = map[string]int64{"Все привычки": 0}
	options := []string{"Все привычки"}
	for _, habit := range habits {
		option := habit.Name
		if _, exists := s.filterIDs[option]; exists {
			option = fmt.Sprintf("%s (#%d)", habit.Name, habit.ID)
		}
		options = append(options, option)
		s.filterIDs[option] = habit.ID
	}
	if _, exists := s.filterIDs[s.filterSelect.Selected]; !exists {
		s.filterSelect.Selected = "Все привычки"
		s.habitFilter = 0
	}
	s.filterSelect.Options = options
	s.filterSelect.Refresh()

	previous := widget.NewButton("<", func() {
		s.month = s.month.AddDate(0, -1, 0)
		s.refreshCalendar()
	})
	next := widget.NewButton(">", func() {
		s.month = s.month.AddDate(0, 1, 0)
		s.refreshCalendar()
	})
	s.nextMonth = next
	if s.month.Before(startOfMonth(startOfToday())) {
		next.Enable()
	} else {
		next.Disable()
	}
	monthLabel := widget.NewLabel(russianMonth(s.month))
	header := container.NewBorder(nil, nil, previous, next, monthLabel)

	weekdays := container.NewGridWithColumns(7,
		widget.NewLabel("Пн"), widget.NewLabel("Вт"), widget.NewLabel("Ср"),
		widget.NewLabel("Чт"), widget.NewLabel("Пт"), widget.NewLabel("Сб"), widget.NewLabel("Вс"),
	)
	cells := make([]fyne.CanvasObject, 0, 42)
	first := startOfMonth(s.month)
	leading := (int(first.Weekday()) + 6) % 7
	daysInMonth := first.AddDate(0, 1, -1).Day()
	for index := 0; index < 42; index++ {
		dayNumber := index - leading + 1
		if dayNumber < 1 || dayNumber > daysInMonth {
			cells = append(cells, widget.NewLabel(""))
			continue
		}
		date := time.Date(first.Year(), first.Month(), dayNumber, 0, 0, 0, 0, first.Location())
		progress, progressErr := s.service.GetDayProgress(context.Background(), date, s.habitFilter)
		if progressErr != nil {
			cells = append(cells, widget.NewLabel("!"))
			continue
		}
		cells = append(cells, newCalendarCell(date, progress, sameDay(date, s.selectedDate), func(value time.Time) {
			if value.After(startOfToday()) {
				return
			}
			s.selectedDate = value
			s.refresh()
		}))
	}
	days := container.NewGridWithColumns(7, cells...)
	legend := widget.NewLabel("Заливка показывает процент выполнения привычек")
	s.calendarRoot.Objects = []fyne.CanvasObject{header, s.filterSelect, weekdays, days, legend}
	s.calendarRoot.Refresh()
}

func (s *HabitScreen) dailyHabitRow(habit storage.DailyHabit) fyne.CanvasObject {
	check := widget.NewCheck(habit.Name, func(checked bool) {
		if s.updating {
			return
		}
		var err error
		if checked {
			err = s.service.SetStatus(context.Background(), s.selectedDate, habit.ID, storage.HabitStatusCompleted)
		} else {
			err = s.service.ClearStatus(context.Background(), s.selectedDate, habit.ID)
		}
		if err != nil {
			dialog.ShowError(err, s.window)
			return
		}
		// Keep the tapped checkbox alive until Fyne finishes dispatching the event.
		fyne.Do(func() {
			s.refresh()
		})
	})
	s.updating = true
	check.SetChecked(habit.Status == storage.HabitStatusCompleted)
	s.updating = false

	if habit.Description != "" {
		description := widget.NewLabel(habit.Description)
		description.Importance = widget.LowImportance
		return container.NewVBox(check, description)
	}
	if habit.Status == storage.HabitStatusSkipped {
		return container.NewVBox(check, widget.NewLabel("Пропущено сегодня"))
	}
	return check
}

func (s *HabitScreen) showHabitForm(habit *storage.Habit) {
	name := widget.NewEntry()
	description := widget.NewMultiLineEntry()
	title := "Новая привычка"
	if habit != nil {
		title = "Изменить привычку"
		name.SetText(habit.Name)
		description.SetText(habit.Description)
	}

	form := dialog.NewForm(
		title,
		"Сохранить",
		"Отмена",
		[]*widget.FormItem{
			widget.NewFormItem("Название", name),
			widget.NewFormItem("Описание", description),
		},
		func(confirmed bool) {
			if !confirmed {
				return
			}
			var err error
			if habit == nil {
				err = s.service.Create(context.Background(), name.Text, description.Text)
				if err == nil {
					// A new habit starts today. Return to today so it is immediately visible.
					s.selectedDate = startOfToday()
				}
			} else {
				err = s.service.Update(context.Background(), habit.ID, name.Text, description.Text)
			}
			if err != nil {
				dialog.ShowError(err, s.window)
				return
			}
			s.refresh()
		},
		s.window,
	)
	form.Show()
}

func (s *HabitScreen) showHabitManager() {
	content := container.NewVBox()
	manager := dialog.NewCustomWithoutButtons("Управление привычками", content, s.window)
	closeButton := widget.NewButton("Закрыть", func() { manager.Hide() })
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(1, 14))
	closeRow := container.NewVBox(spacer, container.NewCenter(closeButton))
	manager.SetButtons([]fyne.CanvasObject{closeRow})

	var refresh func()
	refresh = func() {
		content.Objects = nil
		habits, err := s.service.ListAll(context.Background())
		if err != nil {
			content.Add(widget.NewLabel("Ошибка: " + err.Error()))
			content.Refresh()
			return
		}
		if len(habits) == 0 {
			content.Add(widget.NewLabel("Привычек пока нет."))
		}

		for _, habit := range habits {
			content.Add(s.managementRow(habit, refresh))
		}
		content.Refresh()
	}

	refresh()
	manager.Show()
}

func (s *HabitScreen) managementRow(habit storage.Habit, refresh func()) fyne.CanvasObject {
	status := "активна"
	if habit.ArchivedAt.Valid {
		status = "архивная"
	}
	label := widget.NewLabel(fmt.Sprintf("%s (%s)", habit.Name, status))
	actions := container.NewHBox()
	actions.Add(widget.NewButton("Изменить", func() { s.showHabitForm(&habit) }))
	if habit.ArchivedAt.Valid {
		actions.Add(widget.NewButton("Удалить", func() {
			dialog.ShowConfirm("Удалить привычку?", "История этой привычки будет удалена.", func(ok bool) {
				if !ok {
					return
				}
				if err := s.service.DeleteArchived(context.Background(), habit.ID); err != nil {
					dialog.ShowError(err, s.window)
					return
				}
				refresh()
				s.refresh()
			}, s.window)
		}))
	} else {
		actions.Add(widget.NewButton("Архивировать", func() {
			dialog.ShowConfirm("Архивировать привычку?", "Она перестанет влиять на будущие дни.", func(ok bool) {
				if !ok {
					return
				}
				if err := s.service.Archive(context.Background(), habit.ID); err != nil {
					dialog.ShowError(err, s.window)
					return
				}
				refresh()
				s.refresh()
			}, s.window)
		}))
	}
	return container.NewBorder(nil, nil, label, actions, nil)
}

func startOfToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func startOfMonth(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, value.Location())
}

func sameDay(left, right time.Time) bool {
	return left.Year() == right.Year() && left.YearDay() == right.YearDay()
}

func russianDate(value time.Time) string {
	return value.Format("02.01.2006")
}

func russianMonth(value time.Time) string {
	months := [...]string{
		"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
	}
	return fmt.Sprintf("%s %d", months[value.Month()-1], value.Year())
}
