package workdir

import (
	"bytes"
	"testing"
)

// strange that disabling this makse showVisibleGit pass but enabling it does not
func TestLsGit_nochanges(t *testing.T) {
	wd, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	defer wd.RemoveAll()
	err = wd.Command("git", "add", ".").Run()
	if err != nil {
		t.Fatal("git add failed ", err)
	}
	err = wd.Command("git", "commit", "-m", "hepp").Run()
	if err != nil {
		t.Fatal(err)
	}
	err = wd.LsGit(nil, false)
	if err != nil {
		t.Error(err)
	}
}

func Test_lsGit_error(t *testing.T) {
	wd, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	status, _ := wd.GitStatus()
	err = wd.RemoveAll()
	if err != nil {
		t.Fatal(err)
	}
	err = wd.lsGit(nil, status, false)
	if err == nil {
		t.Fatal(err)
	}
}

func TestLsGit_error(t *testing.T) {
	wd, err := TempDir()
	if err != nil {
		t.Fatal(err)
	}
	err = wd.RemoveAll()
	if err != nil {
		t.Fatal(err)
	}
	err = wd.LsGit(nil, false)
	if err == nil {
		t.Fatal(err)
	}
}

func TestLsGit_colored(t *testing.T) {
	wd, _ := setup()
	defer wd.RemoveAll()
	out := bytes.NewBufferString("\n")
	wd.LsGit(out, true)
	if !bytes.Contains(out.Bytes(), []byte(RED)) {
		t.Fail()
	}
}

func TestLsGit(t *testing.T) {
	wd, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	defer wd.RemoveAll()
	out := bytes.NewBufferString("\n")
	wd.LsGit(out, false)
	exp := `
 M A
   B
   empty/
?? ex/e1
?? ex/e2
MM sub/lev/C
`
	got := out.String()
	if exp != got {
		t.Errorf("Expected \n'%s'\ngot\n'%s'\n", exp, got)
	}
}

func TestGitStatus_err(t *testing.T) {
	wd := WorkDir("/")
	_, err := wd.GitStatus()
	if err == nil {
		t.Error("Expected error from GitStatus when checking non repo")
	}
}

func TestGitStatus_Flags(t *testing.T) {
	tmp, _ := setup()
	defer tmp.RemoveAll()
	status, err := tmp.GitStatus()
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		path, exp string
	}{
		{"A", " M "},
		{"sub/lev/C", "MM "},
		{".hidden", "   "},
	}
	for _, c := range cases {
		got := status.Flags(c.path)
		if c.exp != got {
			t.Errorf("Expected %q, got %q for path %q", c.exp, got, c.path)
		}
	}
}
