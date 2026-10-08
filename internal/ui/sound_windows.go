//go:build windows

package ui

import "syscall"

func playTimerSound() {
	user32 := syscall.NewLazyDLL("user32.dll")
	messageBeep := user32.NewProc("MessageBeep")
	_, _, _ = messageBeep.Call(0xFFFFFFFF)
}
