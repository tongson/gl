//go:build windows

package gl

import (
	"syscall"
)

func setPgid() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}

func sKill(pid int) {
	// noop
}
