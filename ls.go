package dir

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Path struct {
	root string
	w    io.Writer
}

func NewDir() *Path {
	return &Path{".", os.Stdout}
}

func (dir *Path) List() {
	out := make(chan string)
	go func() {
		_ = filepath.Walk(dir.root, visitor(out))
		close(out)
	}()
	for path := range out {
		fmt.Fprintf(dir.w, "%s\n", path)
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
