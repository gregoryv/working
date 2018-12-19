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
	i := bytes.Index(s, append([]byte(path), 0x00))
	if i == -1 {
		return "   "
	}
	return string(s[i-3 : i])
}

func (s GitStatus) IsNew(path string) bool {
	return s.Flags(path) == "?? "
}

func (wd WorkDir) Status() (GitStatus, error) {
	data, err := wd.Command("git", "status", "-z").Output()
	if err != nil {
		return GitStatus([]byte{}), err
	}
	return GitStatus(data), nil
}

func (wd WorkDir) LsGit(w io.Writer, colorize bool) error {
	status, err := wd.Status()
	if err != nil {
		return err
	}
	return wd.lsGit(w, status, colorize)
}

func (wd WorkDir) lsGit(w io.Writer, status GitStatus, colorize bool) error {
	if w == nil {
		w = os.Stdout
	}
	if string(status) == "" {
		return wd.Ls(w)
	}
	visit := showVisibleGit(w, string(wd), status, colorize)
	return filepath.Walk(string(wd), visit)
}

func showVisibleGit(w io.Writer, root string, status GitStatus, colorize bool) filepath.WalkFunc {
	lastDir := ""
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		isHidden := f.Name()[0] == '.'
		isRoot := f.Name() == filepath.Base(root)
		if f.IsDir() && isHidden {
			return filepath.SkipDir
		}
		if isHidden || isRoot {
			return nil
		}

		line := string(path[len(root)+1:])
		if f.IsDir() {
			line += "/"
		}
		isNew := status.IsNew(line)
		flags := status.Flags(line)
		if colorize {
			flags = color(flags)
		}
		if isNew && f.IsDir() {
			fmt.Fprint(w, flags, line, "\n")
			return filepath.SkipDir
		}

		if isNew {
			fmt.Fprint(w, flags, line, "\n")
			return nil
		}

		inSubDir := strings.Contains(line, "/")
		noFlags := status.Flags(line) == "   "
		if !f.IsDir() && inSubDir && noFlags {
			return nil
		}
		if f.IsDir() {
			if lastDir == "" {
				lastDir = line
				return nil
			}
			if strings.Index(line, lastDir) == 0 && flags == "   " {
				return nil
			}
			fmt.Fprint(w, flags, lastDir, "\n")
			lastDir = ""
			return nil
		}
		fmt.Fprint(w, flags, line, "\n")
		return nil
	}
}

const (
	NOCOLOR = "\033[0m"
	RED     = "\033[0;31m"
	GREEN   = "\033[0;32m"
)

func color(flags string) string {
	if flags[0] != '?' {
		return fmt.Sprintf("%s%s%s%s%s ", GREEN, string(flags[0]), RED, string(flags[1]), NOCOLOR)
	}
	return fmt.Sprintf("%s%s%s", RED, flags, NOCOLOR)
}
