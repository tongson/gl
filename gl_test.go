package gl

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRun(T *testing.T) {
	T.Parallel()
	T.Run("gl.Run SimpleOk", func(t *testing.T) {
		var exe RunArg
		exe = RunArg{Exe: "true"}
		if ret, _ := exe.Run(); !ret {
			t.Error("Run() wants `true`")
		}
	})
	T.Run("gl.Run SimpleFail", func(t *testing.T) {
		var exe RunArg
		exe = RunArg{Exe: "false"}
		if ret, _ := exe.Run(); ret {
			t.Error("Run() wants `false`")
		}
	})
	T.Run("gl.Run Timeout Exe Fail", func(t *testing.T) {
		exe := RunArg{
			Exe:     "/bin/sheesh",
			Args:    []string{"-c", "sleep 1"},
			Timeout: 3,
		}
		ret, res := exe.Run()
		if ret {
			t.Error("Run() wants `false`")
		}
		expected := "fork/exec /bin/sheesh: no such file or directory"
		if res.Error != expected {
			t.Errorf("Run() wants `%s`; got `%s`", expected, res.Error)
		}
	})
	T.Run("gl.Run Timeout OK", func(t *testing.T) {
		exe := RunArg{
			Exe:     "/bin/sh",
			Args:    []string{"-c", "sleep 1"},
			Timeout: 3,
		}
		ret, _ := exe.Run()
		if !ret {
			t.Error("Run() wants `true`")
		}
	})
	T.Run("gl.Run Timeout Fail", func(t *testing.T) {
		exe := RunArg{
			Exe:     "/bin/sh",
			Args:    []string{"-c", "sleep 1000"},
			Timeout: 1,
		}
		ret, res := exe.Run()
		if ret {
			t.Error("Run() wants `false`")
		}
		expected := "signal: terminated"
		if res.Error != expected {
			t.Errorf("Run() wants `%s`; got `%s`", expected, res.Error)
		}
	})
	T.Run("gl.Run Dir", func(t *testing.T) {
		exe := RunArg{Exe: "ls", Dir: "/etc"}
		ret, out := exe.Run()
		if !ret {
			t.Error("Run() wants `true`")
		}
		if !strings.Contains(out.Stdout, "passwd") {
			t.Error("Run() output must contain string")
		}
	})
	T.Run("gl.Run Args", func(t *testing.T) {
		args := []string{
			"/etc",
		}
		exe := RunArg{Exe: "ls", Args: args}
		ret, out := exe.Run()
		if !ret {
			t.Error("Run() wants `true`")
		}
		if !strings.Contains(out.Stdout, "passwd") {
			t.Error("Run() output must contain string")
		}
	})
	T.Run("gl.Run Stdout", func(t *testing.T) {
		var sh strings.Builder
		sf := func(o string) {
			sh.WriteString(o)
		}
		args := []string{
			"/etc",
		}
		exe := RunArg{Exe: "ls", Args: args, Stdout: sf}
		ret, _ := exe.Run()
		if !ret {
			t.Error("Run() wants `true`")
		}
		if !strings.Contains(sh.String(), "passwd") {
			t.Error("Run() output must contain string")
		}
	})
	T.Run("gl.Run Stderr", func(t *testing.T) {
		var sh strings.Builder
		sf := func(o string) {
			sh.WriteString(o)
		}
		args := []string{
			"foo", "/etc/shadow",
		}
		exe := RunArg{Exe: "grep", Args: args, Stderr: sf}
		ret, _ := exe.Run()
		if ret {
			t.Error("Run() wants `false`")
		}
		if !strings.Contains(sh.String(), "Permission denied") {
			t.Error("Run() output must contain string")
		}
	})
	T.Run("gl.Run Stdin", func(t *testing.T) {
		var exe RunArg
		input := "foo\n\nbar"
		args := []string{
			"-v",
		}
		exe = RunArg{Exe: "cat", Args: args, Stdin: []byte(input)}
		ret, out := exe.Run()
		if ret != true {
			t.Errorf("Run = %t; want `true`", ret)
		}
		if out.Stdout != "foo\n\nbar" {
			t.Errorf("Run = `%s`; want 'foo\n\nbar'", out.Stdout)
		}
		if out.Stderr != "" {
			t.Errorf("Run = `%s`; want ''", out.Stderr)
		}
		if out.Error != "" {
			t.Errorf("Run = `%s`; want ''", out.Error)
		}
	})
	T.Run("gl.Run Env", func(t *testing.T) {
		env := []string{"FOO=BAR"}
		var exe RunArg
		args := []string{
			"BEGIN{print ENVIRON[\"FOO\"]}",
		}
		exe = RunArg{Exe: "awk", Args: args, Env: env}
		ret, out := exe.Run()
		if ret != true {
			t.Errorf("Run = %t; want `true`", ret)
		}
		if out.Stdout != "BAR\n" {
			t.Errorf("Run = %s; want 'BAR\n'", out.Stdout)
		}
		if out.Stderr != "" {
			t.Errorf("Run = %s; want ''", out.Stderr)
		}
		if out.Error != "" {
			t.Errorf("Run = %s; want ''", out.Error)
		}
	})
	T.Run("gl.Run Stdout integrity", func(t *testing.T) {
		var exe RunArg
		exe = RunArg{Exe: "cat", Args: []string{"utf8.txt"}}
		ret, out := exe.Run()
		if !ret {
			t.Error("Run() wants `true`")
		}
		in := out.Stdout
		sha := RunArg{Exe: "sha256sum", Stdin: []byte(in)}
		_, hash := sha.Run()
		if hash.Stdout != "69a26f4dcd4980e2cc60cc622cbfa64924c48337ec1f34931992166392f94dea  -\n" {
			t.Errorf("Hash did not match. Got: %s", hash.Stdout)
		}
	})
	T.Run("gl.Run Stdout stream integrity", func(t *testing.T) {
		var soStr strings.Builder
		soFn := func(s string) {
			_, _ = soStr.WriteString(s)
		}
		var exe RunArg
		exe = RunArg{Exe: "cat", Args: []string{"utf8.txt"}, Stdout: soFn}
		ret, _ := exe.Run()
		if !ret {
			t.Error("Run() wants `true`")
		}
		sha := RunArg{Exe: "sha256sum", Stdin: []byte(soStr.String())}
		_, hash := sha.Run()
		if hash.Stdout != "69a26f4dcd4980e2cc60cc622cbfa64924c48337ec1f34931992166392f94dea  -\n" {
			t.Errorf("Hash did not match. Got: %s", hash.Stdout)
		}
	})
}

func TestIsFile(t *testing.T) {
	b := IsFile("/etc/passwd")
	if b != true {
		t.Errorf("IsFile() = %t; want `true`", b)
	}
	c := IsFile("/etc/_zzz_")
	if c != false {
		t.Errorf("IsFile() =%t; want `false`", c)
	}
}

func TestIsDir(t *testing.T) {
	b := IsDir("/etc")
	if b != true {
		t.Errorf("IsDir(\"directory\") = %t; want `true`", b)
	}
	c := IsDir("/etc/passwd")
	if c != false {
		t.Errorf("IsDir(\"directory\") = %t; want `false`", b)
	}
}

func TestPipeStr(T *testing.T) {
	T.Parallel()
	T.Run("PipeStr", func(t *testing.T) {
		a := "prefix"
		b := "this"
		c := PipeStr(a, b)
		expected := " prefix │\n" + " prefix │ this\n" + " prefix │"
		if c != expected {
			t.Errorf("Got = `%s`; Did not match expected output: %s", c, expected)
		}
	})
	T.Run("PipeStr Empty", func(t *testing.T) {
		a := "prefix"
		b := ""
		c := PipeStr(a, b)
		expected := " prefix │\n" + " prefix │ \n" + " prefix │"
		if c != expected {
			t.Errorf("Got = `%s`; Did not match expected output: %s", c, expected)
		}
	})
	T.Run("PipeStr Newline", func(t *testing.T) {
		a := "prefix"
		b := "one\n\ntwo"
		c := PipeStr(a, b)
		expected := " prefix │\n" + " prefix │ one\n" + " prefix │ \n" + " prefix │ two\n" + " prefix │"
		if c != expected {
			t.Errorf("Got = `%s`; Did not match expected output: %s", c, expected)
		}
	})
}

func TestInsertStr(t *testing.T) {
	a := "yes"
	b := InsertStr(strings.Split(a, ""), "v", 1)
	r := reflect.DeepEqual(b, strings.Split("yves", ""))
	if r != true {
		t.Errorf("InsertStr() = %s; want '[y v e s]'", b)
	}
}

func TestStringToFile(t *testing.T) {
	e := StringToFile("/dev/null", "foo")
	if e != nil {
		t.Errorf("StringToFile() = %s; want 'nil'", e.Error())
	}
	f := StringToFile("/dev/nil", "foo")
	if f == nil {
		t.Errorf("StringToFile() = %s; want 'error'", f.Error())
	}
}

func TestFileRead(t *testing.T) {
	s := FileRead("/etc/passwd")
	if s == "" {
		t.Error("FileRead() = ''; want 'strings'")
	}
	x := FileRead("/etc/sdfsdf")
	if x != "" {
		t.Error("want ''")
	}
	z := FileRead("/etc/shadow")
	if z != "" {
		t.Error("want ''")
	}
	e := FileRead("/dev/null")
	if e != "" {
		t.Error("want ''")
	}
}

func TestFileLines(t *testing.T) {
	s := FileLines("gl.go")
	if s == nil {
		t.Error("FileLines() = nil; want '[]string'")
	}
	x := FileLines("gl.go")
	if len(x) < 0 {
		t.Error("FileLines() len = 0; want > 0")
	}
	if x[0] != "package gl" {
		t.Error("FileLines(gl.go) did not match")
	}
	if x[1] != "" {
		t.Error("FileLines(gl.go) did not match")
	}
	if x[2] != "import (" {
		t.Error("FileLines(gl.go) did not match")
	}
	y := FileLines("/etc/shadow")
	if len(y) > 1 {
		t.Error("FileLines() > 1; want 0")
	}
	z := FileLines("_zzz_.go")
	if len(z) > 1 {
		t.Error("FileLines() > 1; want 0")
	}
}

func TestFileGlob(t *testing.T) {
	x, _ := FileGlob("./*")
	if x[0] != ".git" {
		t.Error("FileGlob(./*) did not match")
	}
	y, _ := FileGlob(".git/*")
	if y[0] != ".git/COMMIT_EDITMSG" {
		t.Error("FileGlob(.git/*) did not match")
	}
}

func TestPathWalk(t *testing.T) {
	var fs strings.Builder
	var err error
	fnwalk := PathWalker(&fs)
	err = filepath.WalkDir("/etc/environment.d", fnwalk)
	if err != nil {
		t.Error(err.Error())
	}
	err = filepath.WalkDir("/etc", fnwalk)
	if err == nil {
		t.Error("Expected an error here.")
	}
	err = filepath.WalkDir("/dev", fnwalk)
	if err == nil {
		t.Error("Expected an error here.")
	}
	if fs.String() == "" {
		t.Error("PathWalker() = ''; want 'strings'")
	}
}
