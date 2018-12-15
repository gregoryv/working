package dir

import (
	"bytes"
	"io/ioutil"
	"testing"

	. "github.com/gregoryv/qual"
)

func TestDir_List(t *testing.T) {
	If := Wrap(t)

	// setup temporary structure
	tmpPath, err := ioutil.TempDir("", "dirlist")
	If(err != nil).Fatal(err)
	tmp := dir(tmpPath)
	defer tmp.RemoveAll()
	// add files
	fileA, err := tmp.Touch("A")
	If(err != nil).Fatal(err)

	fileB, err := tmp.Touch("B")
	If(err != nil).Fatal(err)

	out := bytes.NewBufferString("")
	d := &Path{root: tmpPath, w: out}
	d.List()
	got := out.String()
	exp := tmp.Join(fileA+"\n") + tmp.Join(fileB+"\n")
	If(exp != got).Errorf("Expected \n%s, got \n%s", exp, got)
}

func TestNewDir(t *testing.T) {
	got := NewDir()
	if got == nil {
		t.Fail()
	}
}
