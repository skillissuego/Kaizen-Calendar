package ui

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type PomodoroScreen struct {
	window       fyne.Window
	workSeconds  int
	breakSeconds int
	remaining    int
	workPhase    bool
	running      bool
	repeat       bool
	soundEnabled bool
	stop         chan struct{}

	phaseLabel   *widget.Label
	timerLabel   *widget.Label
	startButton  *widget.Button
	repeatButton *widget.Button
	applyButton  *widget.Button
	workEntry    *widget.Entry
	breakEntry   *widget.Entry
	root         *fyne.Container
}

func NewPomodoroScreen(window fyne.Window, soundEnabled bool) (fyne.CanvasObject, *PomodoroScreen) {
	screen := &PomodoroScreen{
		window:       window,
		workSeconds:  25 * 60,
		breakSeconds: 5 * 60,
		remaining:    25 * 60,
		workPhase:    true,
		repeat:       true,
		soundEnabled: soundEnabled,
		phaseLabel:   widget.NewLabel("Работа"),
		timerLabel:   widget.NewLabel("25:00"),
		workEntry:    widget.NewEntry(),
		breakEntry:   widget.NewEntry(),
	}
	screen.workEntry.SetText("25")
	screen.breakEntry.SetText("5")

	presets := widget.NewSelect([]string{"25 / 5", "50 / 10", "90 / 15", "Пользовательский"}, func(selected string) {
		switch selected {
		case "25 / 5":
			screen.setPreset(25, 5)
		case "50 / 10":
			screen.setPreset(50, 10)
		case "90 / 15":
			screen.setPreset(90, 15)
		case "Пользовательский":
			screen.workEntry.Enable()
			screen.breakEntry.Enable()
			screen.applyButton.Show()
			screen.workEntry.SetText("")
			screen.breakEntry.SetText("")
		}
	})
	presets.SetSelected("25 / 5")

	applyCustom := widget.NewButton("Применить", func() {
		work, workErr := strconv.Atoi(screen.workEntry.Text)
		breakMinutes, breakErr := strconv.Atoi(screen.breakEntry.Text)
		if workErr != nil || breakErr != nil || work < 1 || breakMinutes < 1 {
			dialog.ShowError(fmt.Errorf("введите положительные целые значения минут"), screen.window)
			return
		}
		screen.setDurations(work, breakMinutes)
	})
	screen.applyButton = applyCustom
	screen.setPreset(25, 5)

	screen.startButton = widget.NewButton("Старт", screen.toggle)
	screen.repeatButton = widget.NewButton("Повтор: вкл", screen.toggleRepeat)
	reset := widget.NewButton("Сбросить", screen.reset)
	controls := container.NewHBox(screen.startButton, screen.repeatButton, reset)
	custom := container.NewBorder(
		nil,
		nil,
		container.NewHBox(widget.NewLabel("Работа"), screen.workEntry, widget.NewLabel("мин")),
		applyCustom,
		container.NewHBox(widget.NewLabel("Перерыв"), screen.breakEntry, widget.NewLabel("мин")),
	)

	content := container.NewVBox(
		presets,
		custom,
		widget.NewSeparator(),
		screen.phaseLabel,
		screen.timerLabel,
		controls,
	)
	screen.root = container.NewCenter(content)
	return screen.root, screen
}

func (s *PomodoroScreen) SetSoundEnabled(enabled bool) {
	s.soundEnabled = enabled
}

func (s *PomodoroScreen) setDurations(workMinutes, breakMinutes int) {
	s.workSeconds = workMinutes * 60
	s.breakSeconds = breakMinutes * 60
	s.reset()
}

func (s *PomodoroScreen) setPreset(workMinutes, breakMinutes int) {
	s.workEntry.SetText(strconv.Itoa(workMinutes))
	s.breakEntry.SetText(strconv.Itoa(breakMinutes))
	s.workEntry.Disable()
	s.breakEntry.Disable()
	if s.applyButton != nil {
		s.applyButton.Hide()
	}
	s.setDurations(workMinutes, breakMinutes)
}

func (s *PomodoroScreen) toggle() {
	if s.running {
		s.running = false
		close(s.stop)
		s.stop = nil
		s.startButton.SetText("Продолжить")
		return
	}
	s.running = true
	s.startButton.SetText("Пауза")
	stop := make(chan struct{})
	s.stop = stop
	go s.run(stop)
}

func (s *PomodoroScreen) toggleRepeat() {
	s.repeat = !s.repeat
	if s.repeat {
		s.repeatButton.SetText("Повтор: вкл")
	} else {
		s.repeatButton.SetText("Повтор: выкл")
	}
}

func (s *PomodoroScreen) run(stop chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fyne.Do(func() {
				if s.running && s.stop == stop {
					s.tick()
				}
			})
		case <-stop:
			return
		}
	}
}

func (s *PomodoroScreen) tick() {
	if s.remaining > 0 {
		s.remaining--
	}
	if s.remaining == 0 {
		if s.soundEnabled {
			playTimerSound()
		}
		if s.repeat {
			s.workPhase = !s.workPhase
			if s.workPhase {
				s.remaining = s.workSeconds
			} else {
				s.remaining = s.breakSeconds
			}
		} else {
			s.stopTimer()
		}
	}
	s.updateDisplay()
}

func (s *PomodoroScreen) stopTimer() {
	if s.running {
		s.running = false
		close(s.stop)
		s.stop = nil
	}
	if s.startButton != nil {
		s.startButton.SetText("Старт")
	}
}

func (s *PomodoroScreen) reset() {
	s.stopTimer()
	s.workPhase = true
	s.remaining = s.workSeconds
	s.updateDisplay()
}

func (s *PomodoroScreen) updateDisplay() {
	minutes := s.remaining / 60
	seconds := s.remaining % 60
	s.timerLabel.SetText(fmt.Sprintf("%02d:%02d", minutes, seconds))
	if s.workPhase {
		s.phaseLabel.SetText("Работа")
	} else {
		s.phaseLabel.SetText("Перерыв")
	}
}
