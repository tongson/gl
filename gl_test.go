package gl

import (
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
		t.Errorf("InsertStr() = %s; want 'yves'", b)
	}
}
