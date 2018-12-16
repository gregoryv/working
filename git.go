package workdir

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type GitStatus []byte

// Flags returns the status letters of a given path if any exist
func (s GitStatus) Flags(path string) string {
	i := bytes.Index(s, []byte(path))
	if i == -1 {
		return ""
	}
	return string(s[i-3 : i])
}

func (wd WorkDir) GitStatus() (GitStatus, error) {
	data, err := wd.Command("git", "status", "-z").Output()
	if err != nil {
		return GitStatus([]byte{}), err
	}
	return GitStatus(data), nil
}

func (wd WorkDir) GitLs(w io.Writer) {
	status, err := wd.GitStatus()
	if err != nil {
		wd.Ls(w)
	}
	visit := showVisibleGit(w, string(wd), status)
	filepath.Walk(string(wd), visit)
}

func showVisibleGit(w io.Writer, root string, status GitStatus) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Index(f.Name(), ".") == 0 {
			if f.IsDir() && path != "." {
				return filepath.SkipDir
			}
			return nil
		}
		if f.Name() != filepath.Base(root) {
			line := string(path[len(root)+1:])
			if f.IsDir() {
				line += "/"
			}
			flags := ""
			if !f.IsDir() {
				flags = status.Flags(line)
			}
			if flags == "" {
				flags = "   "
			}
			fmt.Fprint(w, flags, line, "\n")
		}
		return nil
	}
}
