package main

import (
	"fmt"
	"github.com/tongson/gl"
)

func main() {
  {
		var ret bool
		var stdout string
		var stderr string
		var goerr string
		fmt.Println("[+]", "start gl.RunArgs")
		rargs := gl.RunArgs{
			Exe:     "/bin/sh",
			Args:    []string{"-c", "sleep 6"},
		}
		ret, stdout, stderr, goerr = rargs.Run()
		fmt.Println("ret", ret)
		fmt.Println("stdout", stdout)
		fmt.Println("stderr", stderr)
		fmt.Println("goerr", goerr)
		fmt.Println("[+]", "end gl.RunArgs")
	}
	{
		var ret bool
		var stdout string
		var stderr string
		var goerr string
		fmt.Println("[+]", "start gl.RunArgs Timeout")
		rargs := gl.RunArgs{
			Exe:     "/bin/sh",
			Args:    []string{"-c", "sleep 6"},
			Timeout: 3,
		}
		ret, stdout, stderr, goerr = rargs.Run()
		fmt.Println("ret", ret)
		fmt.Println("stdout", stdout)
		fmt.Println("stderr", stderr)
		fmt.Println("goerr", goerr)
		fmt.Println("[+]", "end gl.RunArgs")
	}
}
