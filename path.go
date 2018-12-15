package dir

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
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
		_ = filepath.Walk(p.root, visitor(out))
		close(out)
	}()
	for path := range out {
		fmt.Fprintf(p.w, "%s\n", path)
	}
}

func visitor(out chan string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		out <- path
		return nil
	}
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
