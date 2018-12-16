package workdir

import (
	"bytes"
	"testing"

	. "github.com/gregoryv/qual"
)

func TestPath_Ls(t *testing.T) {
	tmp, _ := setup()
	defer tmp.RemoveAll()

	out := bytes.NewBufferString("")
	tmp.Writer = out
	tmp.Ls()
	got := out.String()
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
	if got == nil {
		t.Fail()
	}
}

func TestPath_String(t *testing.T) {
	p := New()
	got := p.String()
	exp := "."
	if exp != got {
		t.Errorf("Expected %q, got %q", exp, got)
	}
}
