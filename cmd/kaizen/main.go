package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	kaizenapp "kaizen-calendar/internal/app"
	"kaizen-calendar/internal/storage"
	"kaizen-calendar/internal/ui"
)

//go:embed icon.png
var iconPNG []byte

func main() {
	databasePath, err := databasePath()
	if err != nil {
		showFatalError(err)
		return
	}

	db, err := storage.Open(databasePath)
	if err != nil {
		showFatalError(err)
		return
	}
	defer db.Close()

	habitRepository := storage.NewHabitRepository(db)
	habitService := kaizenapp.NewHabitService(habitRepository)
	todoRepository := storage.NewTodoRepository(db)
	todoService := kaizenapp.NewTodoService(todoRepository)
	settingsRepository := storage.NewSettingsRepository(db)
	settingsService := kaizenapp.NewSettingsService(settingsRepository)
	themeName, _ := settingsService.ThemeName(context.Background())
	soundEnabled, _ := settingsService.SoundEnabled(context.Background())

	kaizenApp := app.NewWithID("com.kaizen.calendar")
	kaizenApp.Settings().SetTheme(ui.NewKaizenTheme(themeName))

	appIcon := fyne.NewStaticResource("icon.png", iconPNG)
	kaizenApp.SetIcon(appIcon)

	window := kaizenApp.NewWindow("Kaizen Calendar")
	window.SetIcon(appIcon)

	pomodoroObject, pomodoroScreen := ui.NewPomodoroScreen(window, soundEnabled)
	settingsScreen := ui.NewSettingsScreen(settingsService, pomodoroScreen, kaizenApp, window, habitService)
	tabs := container.NewAppTabs(
		container.NewTabItem("Kaizen", ui.NewHabitScreen(habitService, window)),
		container.NewTabItem("Pomodoro", pomodoroObject),
		container.NewTabItem("Todo", ui.NewTodoScreen(todoService, window)),
		container.NewTabItem("Настройки", settingsScreen),
	)
	window.SetContent(tabs)
	window.Resize(fyne.NewSize(720, 480))
	window.CenterOnScreen()
	window.ShowAndRun()
}

func databasePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config directory: %w", err)
	}

	appDir := filepath.Join(configDir, "KaizenCalendar")
	if err = os.MkdirAll(appDir, 0700); err != nil {
		return "", fmt.Errorf("create application data directory: %w", err)
	}
	return filepath.Join(appDir, "kaizen.db"), nil
}

func showFatalError(err error) {
	kaizenApp := app.New()
	window := kaizenApp.NewWindow("Kaizen Calendar - ошибка")
	window.SetContent(container.NewCenter(widget.NewLabel(
		"Не удалось запустить приложение:\n\n" + err.Error(),
	)))
	window.Resize(fyne.NewSize(600, 240))
	window.ShowAndRun()
}
