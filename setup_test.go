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
func setup() (wd WorkDir, err error) {
	wd, err = TempDir()
	if err != nil {
		return
	}
	wd.MkdirAll("sub/lev", "empty", "ex")
	_, err = wd.TouchAll("A", "B", "sub/lev/C", ".hidden")
	if err != nil {
		return
	}

	err = wd.Command("git", "init", ".").Run()
	if err != nil {
		return
	}
	wd.Command("git", "add", ".").Run()
	wd.Command("git", "commit", "-m", "Initial").Run()
	wd.WriteFile("A", []byte("hello"))
	wd.WriteFile("sub/lev/C", []byte("hello"))
	wd.Command("git", "add", "sub/lev/C").Run()
	wd.WriteFile("sub/lev/C", []byte("world"))
	wd.TouchAll("ex/e1", "ex/e2")
	return
}

func Test_setup(t *testing.T) {
	wd, err := setup()
	if err != nil {
		t.Error(err)
	}
	defer wd.RemoveAll()
	out, _ := wd.Command("git", "status", "-z").Output()
	got := string(out)
	exp := " M A\x00MM sub/lev/C\x00?? ex/\x00"
	if exp != got {
		t.Errorf("Expected \n%q\ngot\n%q\n", exp, got)
	}
}
