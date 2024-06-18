//go:build windows

package gl

func setPgid() {
	// noop
}

func sKill(pid int) {
	// noop
}
