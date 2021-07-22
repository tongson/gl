package gl

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

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
	fnwalk := PathWalker(&fs)
	filepath.Walk("/usr/bin", fnwalk)
	if fs.String() == "" {
		t.Error("PathWalker() = ''; want 'strings'")
	}
}
