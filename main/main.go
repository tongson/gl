package main

import (
	"github.com/tongson/gl"
)

func main() {
	p := func(s string) {
		print(s)
	}
	script := `
	set -efu
	VAR="XXX"
	export EVAR="YYY"
	sudo --preserve-env=XVAR,VAR -s /bin/dash -c '/usr/bin/env' e
	`
	env := []string{"LC_ALL=C", "XVAR=ZZZ"}
	a := gl.RunArg{Env: env, Dir: "/etc", Exe: "sh", Args: []string{"-c", script}, Stdout: p}
	_, _ = a.Run()
}
