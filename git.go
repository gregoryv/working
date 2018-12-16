package dir

import "strings"

type GitStatus []string

func Parse(body string) GitStatus {
	lines := strings.Split(body, "\n")
	status := GitStatus(lines)
	return status
}
