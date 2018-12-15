package dir

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Path struct {
	root string
	w    io.Writer
}

func NewPath() *Path {
	return &Path{".", os.Stdout}
}

func (p *Path) String() string {
	return p.root
}

func (p *Path) Ls() {
	out := make(chan string)
	go func() {
		_ = filepath.Walk(p.root, p.visitor(out))
		close(out)
	}()
	for path := range out {
		fmt.Fprintf(p.w, "%s\n", path)
	}
}

func (p *Path) visitor(out chan string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip hidden
		if strings.Index(f.Name(), ".") == 0 {
			return nil
		}
		if f.Name() != filepath.Base(p.root) {
			out <- f.Name()
		}
		return nil
	}
}

func (p *Path) TouchAll(filenames ...string) ([]string, error) {
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

func (p *Path) Touch(filename string) (string, error) {
	fh, err := os.Create(path.Join(p.root, filename))
	if err != nil {
		return filename, err
	}
	return filename, fh.Close()
}

func (p *Path) Join(filename string) string {
	return filepath.Join(p.root, filename)
}

func (p *Path) RemoveAll() error {
	return os.RemoveAll(p.root)
}
