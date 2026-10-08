//go:build linux

package app

import (
	"fmt"
	"os"
	"path/filepath"
)

const autostartEntryName = "kaisen-calendar.desktop"

// autostartEnabled reports whether the XDG autostart entry for Kaizen exists.
func autostartEnabled() bool {
	path, err := autostartEntryPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

// setAutostart enables or disables autostart by writing or removing an
// XDG .desktop entry in the user's autostart directory.
func setAutostart(enabled bool) error {
	path, err := autostartEntryPath()
	if err != nil {
		return err
	}

	if !enabled {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("disable Linux autostart: %w", err)
		}
		return nil
	}

	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}

	content := "[Desktop Entry]\n" +
		"Type=Application\n" +
		"Name=Kaizen Calendar\n" +
		"Comment=Календарь привычек, Pomodoro и задачи\n" +
		"Exec=" + executable + "\n" +
		"Terminal=false\n" +
		"X-GNOME-Autostart-enabled=true\n"

	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return fmt.Errorf("enable Linux autostart: %w", err)
	}
	return nil
}

// autostartEntryPath returns the path to the Kaizen .desktop autostart entry.
func autostartEntryPath() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get user home directory: %w", err)
		}
		configHome = filepath.Join(home, ".config")
	}

	dir := filepath.Join(configHome, "autostart")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("create autostart directory: %w", err)
	}
	return filepath.Join(dir, autostartEntryName), nil
}