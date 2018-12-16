package dir

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestParse(t *testing.T) {
	status := bytes.NewBufferString(``)
	_ = Parse(status.String())
}

// setup creates a repository with some files
//
// empty/
// sub/
// sub/C     <--- Modified with cache
// A         <--- Modified
// B
// .hidden
func setup() (tmpPath string, err error) {
	tmpPath, err = ioutil.TempDir("", "gitstatus")
	if err != nil {
		return
	}
	tmp := &Path{Root: tmpPath, w: &NopWriter{}}
	//defer tmp.RemoveAll()

	tmp.MkdirAll("sub", "empty")
	_, err = tmp.TouchAll("A", "B", "sub/C", ".hidden")
	if err != nil {
		return
	}
	os.Chdir(tmpPath)
	err = exec.Command("git", "init", ".").Run()
	if err != nil {
		return
	}
	exec.Command("git", "add", ".").Run()
	exec.Command("git", "commit", "-m", "Initial").Run()
	tmp.WriteFile("A", []byte("hello"))
	tmp.WriteFile("sub/C", []byte("hello"))
	exec.Command("git", "add", "sub/C").Run()
	tmp.WriteFile("sub/C", []byte("world"))
	return
}

func Test_setup(t *testing.T) {
	tmpPath, err := setup()
	if err != nil {
		fmt.Println(err)
	}
	os.Chdir(tmpPath)
	out, _ := exec.Command("git", "status", "-z").Output()
	got := string(out)
	exp := " M A\x00MM sub/C\x00"
	if exp != got {
		t.Errorf("Expected \n%q\ngot\n%q\n", exp, got)
	}
}
