//go:build linux

package ui

import "os/exec"

// playTimerSound plays the freedesktop alarm sound using the first
// available system player, so the Pomodoro timer alert works on most
// Linux desktops without extra dependencies.
//
// Direct file players (paplay/aplay) are preferred over canberra-gtk-play
// because the latter relies on desktop sound-theme events, which may be
// disabled and silently drop the alert.
func playTimerSound() {
	commands := [][]string{
		{"paplay", "/usr/share/sounds/freedesktop/stereo/alarm-clock-elapsed.oga"},
		{"aplay", "/usr/share/sounds/alsa/Front_Center.wav"},
		{"canberra-gtk-play", "-i", "alarm-clock-elapsed"},
	}
	for _, args := range commands {
		if _, err := exec.LookPath(args[0]); err != nil {
			continue
		}
		cmd := exec.Command(args[0], args[1:]...)
		if err := cmd.Start(); err != nil {
			continue
		}
		_ = cmd.Process.Release()
		return
	}
}
