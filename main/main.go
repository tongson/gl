package main

import (
	"github.com/tongson/gl"
)

func main() {
	p := func(s string) {
		print(s)
	}
	a := gl.RunArg{Dir: "/etc", Exe: "ls", Stdout: p}
	_, _ = a.Run()
}
