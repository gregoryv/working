package workdir

import (
	"bytes"
	"testing"
)

func Test_LsGit_error(t *testing.T) {
	wd, _ := TempDir()
	wd.RemoveAll()
	err := wd.LsGit(nil, false)
	if err == nil {
		t.Fail()
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
	tmp, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	defer tmp.RemoveAll()
	out := bytes.NewBufferString("\n")
	tmp.LsGit(out, false)
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

// hlh ljlk
