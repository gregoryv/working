package workdir

import (
	"bytes"
	"testing"
)

func TestGitLs(t *testing.T) {
	tmp, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	defer tmp.RemoveAll()
	out := bytes.NewBufferString("\n")
	tmp.Writer = out
	tmp.GitLs()
	exp := `
 M A
   B
   empty/
   sub/
MM sub/C
`
	got := out.String()
	if exp != got {
		t.Errorf("Expected \n'%s'\ngot\n'%s'\n", exp, got)
	}
}

func TestGitStatus_err(t *testing.T) {
	wd := New()
	wd.Root = "/"
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
		{"sub/C", "MM "},
		{".hidden", ""},
	}
	for _, c := range cases {
		got := status.Flags(c.path)
		if c.exp != got {
			t.Errorf("Expected %q, got %q for path %q", c.exp, got, c.path)
		}
	}
}
