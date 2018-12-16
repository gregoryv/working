package workdir

import (
	"testing"
)

// setup creates a repository with some files
//
// empty/
// sub/
// sub/C     <--- Modified with cache
// A         <--- Modified
// B
// .hidden
func setup() (tmp *WorkDir, err error) {
	tmp, err = TempDir()
	if err != nil {
		return
	}
	tmp.MkdirAll("sub", "empty")
	_, err = tmp.TouchAll("A", "B", "sub/C", ".hidden")
	if err != nil {
		return
	}

	err = tmp.Command("git", "init", ".").Run()
	if err != nil {
		return
	}
	tmp.Command("git", "add", ".").Run()
	tmp.Command("git", "commit", "-m", "Initial").Run()
	tmp.WriteFile("A", []byte("hello"))
	tmp.WriteFile("sub/C", []byte("hello"))
	tmp.Command("git", "add", "sub/C").Run()
	tmp.WriteFile("sub/C", []byte("world"))
	return
}

func Test_setup(t *testing.T) {
	tmp, err := setup()
	if err != nil {
		t.Error(err)
	}
	out, _ := tmp.Command("git", "status", "-z").Output()
	got := string(out)
	exp := " M A\x00MM sub/C\x00"
	if exp != got {
		t.Errorf("Expected \n%q\ngot\n%q\n", exp, got)
	}
}
