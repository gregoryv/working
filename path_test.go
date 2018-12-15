package dir

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/gregoryv/qual"
)

func TestPath_Ls(t *testing.T) {
	If := Wrap(t)

	// setup temporary structure
	tmpPath, err := ioutil.TempDir("", "dirlist")
	If(err != nil).Fatal(err)

	tmp := &Path{root: tmpPath, w: &NopWriter{}}
	defer tmp.RemoveAll()

	files, err := tmp.TouchAll("A", "B", ".hidden")
	If(err != nil).Fatal(err)
	// end setup

	out := bytes.NewBufferString("")
	d := &Path{root: tmpPath, w: out}
	d.Ls()
	got := out.String()
	exp := strings.Join(files[:len(files)-1], "\n") + "\n"
	If(exp != got).Errorf("Expected \n'%s'\ngot \n'%s'", exp, got)
}

func TestNewPath(t *testing.T) {
	got := NewPath()
	if got == nil {
		t.Fail()
	}
}

type NopWriter struct{}

func (nw *NopWriter) Write(p []byte) (n int, err error) {
	return
}
