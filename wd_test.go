package workdir

import (
	"bytes"
	"os"
	"testing"

	. "github.com/gregoryv/qual"
)

func TestTempDir_error(t *testing.T) {
	os.Setenv("TMPDIR", "/_no_such_dir")
	defer os.Setenv("TMPDIR", "/tmp")
	_, err := TempDir()
	if err == nil {
		t.Fail()
	}
}

func TestPath_Ls(t *testing.T) {
	wd, _ := setup()
	defer wd.RemoveAll()

	w := bytes.NewBufferString("")
	wd.Ls(w)
	got := w.String()
	exp := `A
B
empty/
ex/
sub/
`
	If := Wrap(t)
	If(exp != got).Errorf("Expected \n'%s'\ngot \n'%s'", exp, got)
}

func TestNew(t *testing.T) {
	got := New(".")
	if got != "." {
		t.Fail()
	}
}

func TestPath_String(t *testing.T) {
	got := WorkDir(".").String()
	exp := "."
	if exp != got {
		t.Errorf("Expected %q, got %q", exp, got)
	}
}

func TestTouch(t *testing.T) {
	wd, _ := TempDir()
	wd.MkdirAll("x")
	os.Chmod(wd.Join("x"), 0000)
	_, err := wd.TouchAll("x")
	if err == nil {
		t.Error("Expected to fail")
	}
	// Cleanup
	os.Chmod(wd.Join("x"), 0644)
	wd.RemoveAll()
}

func TestMkdirAll(t *testing.T) {
	wd, _ := TempDir()
	os.Chmod(string(wd), 0000)
	err := wd.MkdirAll("hepp")
	if err == nil {
		t.Error("Expected to fail")
	}
	os.Chmod(wd.Join("x"), 0644)
	wd.RemoveAll()
}
