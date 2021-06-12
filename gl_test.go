package gl

import (
	"reflect"
	"path/filepath"
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
		t.Errorf("FileRead() = %s; want 'strings'", s)
	}
}

func TestPathWalk(t *testing.T) {
	var fs strings.Builder
	fnwalk := PathWalker(&fs)
	filepath.Walk("/etc/sysctl.d", fnwalk)
	if fs.String() == "" {
		t.Error("PathWalker() = ''; want 'strings'")
	}
}
