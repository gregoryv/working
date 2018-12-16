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

type WorkDir struct {
	Root   string
	Writer io.Writer
	Skip   func(string, os.FileInfo) bool
	Format func(string, os.FileInfo) string
	Filter func(string) string // Filters final output
}

func Hidden(path string, f os.FileInfo) bool {
	if strings.Index(f.Name(), ".") == 0 {
		return true
	}
	return false
}

func unfiltered(path string) string              { return path }
func nameOnly(path string, f os.FileInfo) string { return f.Name() }

func New() *WorkDir {
	return &WorkDir{Root: ".", Writer: os.Stdout, Skip: Hidden, Format: nameOnly,
		Filter: unfiltered}
}

// Returns a new temporary working directory.
func TempDir() (wd *WorkDir, err error) {
	tmpPath, err := ioutil.TempDir("", "workdir")
	if err != nil {
		return
	}
	wd = New()
	wd.Root = tmpPath
	wd.Writer = &NopWriter{}
	return
}

func (p *WorkDir) WriteFile(file string, data []byte) error {
	return ioutil.WriteFile(p.Join(file), data, 0644)
}

func (p *WorkDir) MkdirAll(subDirs ...string) error {
	for _, sub := range subDirs {
		err := os.MkdirAll(filepath.Join(p.Root, sub), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *WorkDir) Command(cmd string, args ...string) *exec.Cmd {
	os.Chdir(p.Root)
	return exec.Command(cmd, args...)
}

func (p *WorkDir) String() string {
	return p.Root
}

func (p *WorkDir) Ls() {
	out := make(chan string)
	go func() {
		_ = filepath.Walk(p.Root, p.visitor(out))
		close(out)
	}()
	for path := range out {
		line := p.Filter(path)
		if line != "" {
			fmt.Fprintf(p.Writer, "%s\n", line)
		}
	}
}

func (p *WorkDir) visitor(out chan string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if p.Skip(path, f) {
			if f.IsDir() && path != "." {
				return filepath.SkipDir
			}
			return nil
		}
		if f.Name() != filepath.Base(p.Root) {
			out <- p.Format(path, f)
		}
		return nil
	}
}

func (p *WorkDir) TouchAll(filenames ...string) ([]string, error) {
	files := make([]string, len(filenames))
	for i, name := range filenames {
		name, err := p.Touch(name)
		if err != nil {
			return files, err
		}
		files[i] = name
	}
	return files, nil
}

func (p *WorkDir) Touch(filename string) (string, error) {
	fh, err := os.Create(path.Join(p.Root, filename))
	if err != nil {
		return filename, err
	}
	return filename, fh.Close()
}

func (p *WorkDir) Join(filename string) string {
	return filepath.Join(p.Root, filename)
}

func (p *WorkDir) RemoveAll() error {
	return os.RemoveAll(p.Root)
}
