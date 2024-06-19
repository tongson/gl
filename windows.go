//go:build windows

package gl

import (
	"syscall"
)

func setPgid() *syscall.SysProcAttr {
	// noop
	return &syscall.SysProcAttr{}
}

func killPgid(pid int) {
	// noop
}
