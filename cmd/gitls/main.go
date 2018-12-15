package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gregoryv/dir"
)

func main() {
	wd := dir.NewPath()

	status, err := exec.Command("git", "-C", wd.Root, "status", "-s").Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	statLines := bytes.Split(status, []byte("\n"))
	wd.Format = func(path string, f os.FileInfo) string {
		line := path
		if f.IsDir() {
			return "   " + line + "/"
		}
		if l := findLine(statLines, []byte(path)); l != path {
			return l
		}
		return "   " + line
	}

	var last string
	wd.Filter = func(line string) string {
		// same folder as before
		sameDir := len(last) != 0 && strings.Index(line, last) == 3
		isDir := line[len(line)-1] == '/'
		if sameDir && isDir {
			return ""
		}
		if last == "" && isDir {
			last = line[3:len(line)]
			return ""
		}
		return line

	}
	wd.Ls()
}

func findLine(lines [][]byte, path []byte) string {
	for _, line := range lines {
		i := bytes.Index(line, path)
		if i > -1 {
			return string(line)
		}
	}
	return string(path)
}
