//go:build windows

package gl

func setPgid() bool {
	// noop
	return nil
}

func sKill(pid int) {
	// noop
}
