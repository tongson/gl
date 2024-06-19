//go:build !windows

package gl

import (
	"syscall"
)

func setPgid() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}

func killPgid(pid int) {
	_ = syscall.Kill(-pid, syscall.SIGTERM)
}
