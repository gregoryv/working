package workdir

import "bytes"

type GitStatus []byte

func (s GitStatus) Flags(path string) string {
	i := bytes.Index(s, []byte(path))
	if i == -1 {
		return ""
	}
	return string(s[i-3 : i])
}
