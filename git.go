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

func (wd WorkDir) GitStatus() (GitStatus, error) {
	data, err := wd.Command("git", "status", "-z").Output()
	if err != nil {
		return GitStatus([]byte{}), err
	}
	return GitStatus(data), nil
}

func (wd WorkDir) LsGit(w io.Writer, colorize bool) {
	if w == nil {
		w = os.Stdout
	}
	status, err := wd.GitStatus()
	if err != nil || string(status) == "" {
		wd.Ls(w)
		return
	}
	visit := showVisibleGit(w, string(wd), status, colorize)
	filepath.Walk(string(wd), visit)
}

func showVisibleGit(w io.Writer, root string, status GitStatus, colorize bool) filepath.WalkFunc {
	lastDir := ""
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
		if f.Name() == filepath.Base(root) {
			return nil
		}

		line := string(path[len(root)+1:])
		flags := status.Flags(line)
		if colorize {
			flags = color(flags)
		}
		if f.IsDir() {
			line += "/"
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
	return fmt.Sprintf("%s%s%s%s%s ", GREEN, string(flags[0]), RED, string(flags[1]), NOCOLOR)
}
