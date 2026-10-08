//go:build windows

package app

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func autostartEnabled() bool {
	command := exec.Command("reg.exe", "query", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, "/v", "KaizenCalendar")
	return command.Run() == nil
}

func setAutostart(enabled bool) error {
	if enabled {
		executable, err := os.Executable()
		if err != nil {
			return fmt.Errorf("get executable path: %w", err)
		}
		path := strings.ReplaceAll(executable, `"`, `\"`)
		command := exec.Command(
			"reg.exe", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`,
			"/v", "KaizenCalendar", "/t", "REG_SZ", "/d", `"`+path+`"`, "/f",
		)
		if output, err := command.CombinedOutput(); err != nil {
			return fmt.Errorf("enable Windows autostart: %w (%s)", err, strings.TrimSpace(string(output)))
		}
		return nil
	}

	command := exec.Command(
		"reg.exe", "delete", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`,
		"/v", "KaizenCalendar", "/f",
	)
	if output, err := command.CombinedOutput(); err != nil {
		return fmt.Errorf("disable Windows autostart: %w (%s)", err, strings.TrimSpace(string(output)))
	}
	return nil
}