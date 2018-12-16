package workdir

import (
	"testing"
)

func TestGitStatus_Flags(t *testing.T) {
	tmp, _ := setup()
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
