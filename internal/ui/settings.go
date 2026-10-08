package ui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	kaizenapp "kaizen-calendar/internal/app"
)

type SettingsScreen struct {
	service     *kaizenapp.SettingsService
	application fyne.App
	window      fyne.Window
	pomodoro    *PomodoroScreen
	statistics  fyne.CanvasObject
	root        fyne.CanvasObject
}

func NewSettingsScreen(service *kaizenapp.SettingsService, pomodoro *PomodoroScreen, application fyne.App, window fyne.Window, habitService *kaizenapp.HabitService) fyne.CanvasObject {
	screen := &SettingsScreen{service: service, window: window, pomodoro: pomodoro, application: application}
	soundEnabled, err := service.SoundEnabled(context.Background())
	if err != nil {
		soundEnabled = true
	}

	sound := widget.NewCheck("Звуковой сигнал Pomodoro", func(enabled bool) {
		if err := screen.service.SetSoundEnabled(context.Background(), enabled); err != nil {
			dialog.ShowError(err, screen.window)
			return
		}
		screen.pomodoro.SetSoundEnabled(enabled)
	})
	sound.SetChecked(soundEnabled)

	themeNames := make([]string, 0, len(ThemeOptions()))
	themeByTitle := make(map[string]string)
	for _, option := range ThemeOptions() {
		themeNames = append(themeNames, option.Title)
		themeByTitle[option.Title] = option.Name
	}
	themeName, err := service.ThemeName(context.Background())
	if err != nil {
		themeName = ThemeTurquoiseDark
	}
	themeSelect := widget.NewSelect(themeNames, func(selected string) {
		name := themeByTitle[selected]
		if err := screen.service.SetThemeName(context.Background(), name); err != nil {
			dialog.ShowError(err, screen.window)
			return
		}
		screen.application.Settings().SetTheme(NewKaizenTheme(name))
	})
	themeSelect.SetSelected(ThemeTitle(themeName))

	autostart := widget.NewCheck("Запускать приложение вместе с системой", func(enabled bool) {
		if err := screen.service.SetAutostartEnabled(enabled); err != nil {
			dialog.ShowError(err, screen.window)
			return
		}
	})
	autostart.SetChecked(service.AutostartEnabled())

	if habitService != nil {
		screen.statistics = NewStatisticsScreen(habitService)
	}

	help := widget.NewLabel("Настройки сохраняются локально на этом компьютере.")
	settingsContent := container.NewVBox(
		widget.NewLabel("Настройки"),
		container.NewBorder(nil, nil, widget.NewLabel("Тема интерфейса:"), nil, themeSelect),
		sound,
		autostart,
		help,
	)
	if screen.statistics != nil {
		settingsContent.Add(widget.NewSeparator())
		settingsContent.Add(widget.NewLabel("Статистика по привычкам"))
		settingsContent.Add(screen.statistics)
	}
	screen.root = container.NewVScroll(settingsContent)
	return screen.root
}

func (s *SettingsScreen) Refresh() {
	// Settings are written immediately, so there is no cached view to reload.
}
