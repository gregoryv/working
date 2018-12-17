package workdir

import (
	"bytes"
	"testing"

	. "github.com/gregoryv/qual"
)

func TestPath_Ls(t *testing.T) {
	wd, _ := setup()
	defer wd.RemoveAll()

	w := bytes.NewBufferString("")
	wd.Ls(w)
	got := w.String()
	exp := `A
B
empty/
sub/
sub/C
`
	If := Wrap(t)
	If(exp != got).Errorf("Expected \n'%s'\ngot \n'%s'", exp, got)
}

func TestNew(t *testing.T) {
	got := New()
	if got != "." {
		t.Fail()
	}
}

func TestPath_String(t *testing.T) {
	got := New().String()
	exp := "."
	if exp != got {
		t.Errorf("Expected %q, got %q", exp, got)
	}
}
