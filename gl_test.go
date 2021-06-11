package gl

import (
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

