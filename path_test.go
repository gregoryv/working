package dir

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/gregoryv/qual"
)

func TestPath_Ls(t *testing.T) {
	If := Wrap(t)

	// setup temporary structure
	tmpPath, err := ioutil.TempDir("", "dirlist")
	subDir := "sub"
	err = os.MkdirAll(filepath.Join(tmpPath, subDir), 0644)
	If(err != nil).Fatal(err)

	tmp := &Path{root: tmpPath, w: &NopWriter{}}
	defer tmp.RemoveAll()

	files, err := tmp.TouchAll("A", "B")
	If(err != nil).Fatal(err)
	files = append(files, subDir)
	tmp.Touch(".hidden")
	// end setup

	out := bytes.NewBufferString("")
	d := &Path{root: tmpPath, w: out, form: nameOnly}
	d.Ls()
	got := out.String()
	exp := strings.Join(files, "\n") + "\n"
	If(exp != got).Errorf("Expected \n'%s'\ngot \n'%s'", exp, got)
}

func TestNewPath(t *testing.T) {
	got := NewPath()
	if got == nil {
		t.Fail()
	}
}

func TestPath_String(t *testing.T) {
	p := NewPath()
	got := p.String()
	exp := "."
	if exp != got {
		t.Errorf("Expected %q, got %q", exp, got)
	}
}

type NopWriter struct{}

func (nw *NopWriter) Write(p []byte) (n int, err error) {
	return
}