//go:build !windows && !linux

package app

import "fmt"

func autostartEnabled() bool {
	return false
}

func setAutostart(enabled bool) error {
	return fmt.Errorf("автозапуск не поддерживается на этой платформе")
}