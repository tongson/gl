package gl

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRunSimpleOk(t *testing.T) {
	var exe RunArgs
	exe = RunArgs{Exe: "true"}
	if ret, _ := exe.Run(); !ret {
		t.Errorf("Run() wants `true`")
	}
}

func TestRunSimpleFail(t *testing.T) {
	var exe RunArgs
	exe = RunArgs{Exe: "false"}
	if ret, _ := exe.Run(); ret {
		t.Errorf("Run() wants `false`")
	}
}

func TestRunStdin(t *testing.T) {
	var exe RunArgs
	input := "foo\n\nbar"
	args := []string{
		"-v",
	}
	exe = RunArgs{Exe: "cat", Args: args, Stdin: []byte(input)}
	ret, out := exe.Run()
	if ret != true {
		t.Errorf("Run = %t; want `true`", ret)
	}
	if out.Stdout != "foo\n\nbar" {
		t.Errorf("Run = %s; want 'foo\n\nbar'", out.Stdout)
	}
	if out.Stderr != "" {
		t.Errorf("Run = %s; want ''", out.Stderr)
	}
	if out.Err != "" {
		t.Errorf("Run = %s; want ''", out.Err)
	}
}

func TestRunEnv(t *testing.T) {
	env := []string{"FOO=BAR"}
	var exe RunArgs
	args := []string{
		"BEGIN{print ENVIRON[\"FOO\"]}",
	}
	exe = RunArgs{Exe: "awk", Args: args, Env: env}
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
	if out.Err != "" {
		t.Errorf("Run = %s; want ''", out.Err)
	}
}

func TestIsFile(t *testing.T) {
	b := IsFile("/etc/passwd")
	if b != true {
		t.Errorf("IsFile(\"\") = %t; want `true`", b)
	}
}

func TestIsDir(t *testing.T) {
	b := IsDir("/etc")
	if b != true {
		t.Errorf("IsDir(\"directory\") = %t; want `true`", b)
	}
}

func TestStatPathFile(t *testing.T) {
	is_file := StatPath("")
	b := is_file("/etc/passwd")
	if b != true {
		t.Errorf("StatPath(\"\") = %t; want `true`", b)
	}
}

func TestStatPathDir(t *testing.T) {
	is_dir := StatPath("directory")
	b := is_dir("/etc")
	if b != true {
		t.Errorf("StatPath(\"directory\") = %t; want `true`", b)
	}
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
}

func TestFileGlob(t *testing.T) {
	x, _ := FileGlob("./*")
	if x[0] != ".git" {
		t.Error("FileGlob(./*) did not match")
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
	if fs.String() == "" {
		t.Error("PathWalker() = ''; want 'strings'")
	}
}
