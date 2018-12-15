package dir

import (
	"bytes"
	"io/ioutil"
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
	// add files
	fileA, err := tmp.Touch("A")
	If(err != nil).Fatal(err)

	fileB, err := tmp.Touch("B")
	If(err != nil).Fatal(err)

	out := bytes.NewBufferString("")
	d := &Path{root: tmpPath, w: out}
	d.Ls()
	got := out.String()
	exp := tmp.Join(fileA+"\n") + tmp.Join(fileB+"\n")
	If(exp != got).Errorf("Expected \n%s, got \n%s", exp, got)
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
