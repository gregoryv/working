package dir

import (
	"testing"
)

func TestGitStatus_Flags(t *testing.T) {
	tmp, _ := setup()
	data, _ := tmp.Command("git", "status", "-z").Output()
	status := GitStatus(data)
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
