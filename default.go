//go:build !windows

package gl

import (
	"syscall"
)

func setPgid() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}

func sKill(pid int) {
	_ = syscall.Kill(-pid, syscall.SIGTERM)
}
