package dir

import (
	"os"
	"path"
	"path/filepath"
)

type dir string

func (d dir) Touch(filename string) (string, error) {
	fh, err := os.Create(path.Join(string(d), filename))
	if err != nil {
		return filename, err
	}
	return filename, fh.Close()
}

func (d dir) Join(filename string) string {
	return filepath.Join(string(d), filename)
}

func (d dir) RemoveAll() error {
	return os.RemoveAll(string(d))
}
