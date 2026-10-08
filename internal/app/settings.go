package app

import (
	"context"
	"strconv"

	"kaizen-calendar/internal/storage"
)

const pomodoroSoundKey = "pomodoro_sound_enabled"
const themeKey = "theme"

type SettingsService struct {
	repository *storage.SettingsRepository
}

func NewSettingsService(repository *storage.SettingsRepository) *SettingsService {
	return &SettingsService{repository: repository}
}

func (s *SettingsService) SoundEnabled(ctx context.Context) (bool, error) {
	value, err := s.repository.Get(ctx, pomodoroSoundKey, "true")
	if err != nil {
		return false, err
	}
	enabled, err := strconv.ParseBool(value)
	if err != nil {
		return true, nil
	}
	return enabled, nil
}

func (s *SettingsService) SetSoundEnabled(ctx context.Context, enabled bool) error {
	return s.repository.Set(ctx, pomodoroSoundKey, strconv.FormatBool(enabled))
}

func (s *SettingsService) ThemeName(ctx context.Context) (string, error) {
	return s.repository.Get(ctx, themeKey, "turquoise-dark")
}

func (s *SettingsService) SetThemeName(ctx context.Context, name string) error {
	return s.repository.Set(ctx, themeKey, name)
}

func (s *SettingsService) AutostartEnabled() bool {
	return autostartEnabled()
}

func (s *SettingsService) SetAutostartEnabled(enabled bool) error {
	return setAutostart(enabled)
}
