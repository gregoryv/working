package workdir

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type WorkDir string

func (wd WorkDir) Ls(w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	return filepath.Walk(wd.String(), showVisible(w, string(wd)))
}

func showVisible(w io.Writer, root string) filepath.WalkFunc {
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
				fmt.Fprint(w, line, "/\n")
				return filepath.SkipDir
			}
			fmt.Fprint(w, line, "\n")
		}
		return nil
	}
}

func New(path string) WorkDir {
	return WorkDir(path)
}

// Returns a new temporary working directory.
func TempDir() (WorkDir, error) {
	tmpPath, err := ioutil.TempDir("", "workdir")
	if err != nil {
		return WorkDir(""), err
	}
	return WorkDir(tmpPath), nil
}

func (wd WorkDir) WriteFile(file string, data []byte) error {
	return ioutil.WriteFile(wd.Join(file), data, 0644)
}

func (wd WorkDir) MkdirAll(subDirs ...string) error {
	for _, sub := range subDirs {
		err := os.MkdirAll(filepath.Join(wd.String(), sub), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wd WorkDir) Command(cmd string, args ...string) *exec.Cmd {
	os.Chdir(wd.String())
	return exec.Command(cmd, args...)
}

func (wd WorkDir) String() string {
	return string(wd)
}

func (wd WorkDir) TouchAll(filenames ...string) ([]string, error) {
	files := make([]string, len(filenames))
	for i, name := range filenames {
		name, err := wd.Touch(name)
		if err != nil {
			return files, err
		}
		files[i] = name
	}
	return files, nil
}

func (wd WorkDir) Touch(filename string) (string, error) {
	fh, err := os.Create(path.Join(wd.String(), filename))
	if err != nil {
		return filename, err
	}
	return filename, fh.Close()
}

func (wd WorkDir) Join(filename string) string {
	return filepath.Join(wd.String(), filename)
}

func (wd WorkDir) RemoveAll() error {
	return os.RemoveAll(wd.String())
}
