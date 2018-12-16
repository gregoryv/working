package workdir

import (
	"bytes"
	"io"
)

type GitStatus []byte

func (s GitStatus) Flags(path string) string {
	i := bytes.Index(s, []byte(path))
	if i == -1 {
		return ""
	}
	return string(s[i-3 : i])
}

func (wd *WorkDir) GitStatus() (GitStatus, error) {
	data, err := wd.Command("git", "status", "-z").Output()
	if err != nil {
		return GitStatus([]byte{}), err
	}
	return GitStatus(data), nil
}

func (wd *WorkDir) GitLs(w io.Writer) {
	wd.Ls(w)
}
