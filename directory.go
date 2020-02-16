package working

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Directory struct {
	path string
}

func (wd *Directory) Chmod(filename string, mode os.FileMode) error {
	return os.Chmod(wd.Join(filename), mode)
}

// IsEmpty returns true if the dir is empty, false otherwise.
// Use empty string to check the workdir itself.
func (wd *Directory) IsEmpty(dir string) bool {
	name := wd.Join(dir)
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false
}

// List content using the given writer. If w is nil stdout is used.
func (wd *Directory) Ls(w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	path := wd.Path()
	return filepath.Walk(path, showVisible(w, path))
}

func showVisible(w io.Writer, root string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		ok, err := isVisibleAndOk(root, path, f)
		if !ok {
			return err
		}

		line := string(path[len(root)+1:])
		if f.IsDir() {
			fmt.Fprint(w, line, "/\n")
			return filepath.SkipDir
		}
		fmt.Fprint(w, line, "\n")
		return nil
	}
}

func isVisibleAndOk(root, path string, f os.FileInfo) (ok bool, err error) {
	if f.Name() == filepath.Base(root) {
		return false, nil
	}
	if !visible(f) {
		if f.IsDir() {
			return false, filepath.SkipDir
		}
		return false, nil
	}
	return true, nil
}

func visible(f os.FileInfo) bool {
	return strings.Index(f.Name(), ".") != 0
}

// Returns a new temporary working directory.
func TempDir() (*Directory, error) {
	tmpPath, err := ioutil.TempDir("", "workdir")
	dir := new(Directory)
	if err != nil {
		return dir, err
	}
	dir.path = tmpPath
	return dir, nil
}

// WriteFile creates/writes over the file with mode 0644
func (wd *Directory) WriteFile(file string, data []byte) error {
	return ioutil.WriteFile(wd.Join(file), data, 0644)
}

// ReadAll loads the given file like ioutil.ReadAll
func (wd *Directory) Load(file string) ([]byte, error) {
	fh, err := os.Open(wd.Join(file))
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	return ioutil.ReadAll(fh)
}

func (wd *Directory) MkdirAll(subDirs ...string) error {
	for _, sub := range subDirs {
		err := os.MkdirAll(filepath.Join(wd.Path(), sub), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wd *Directory) String() string {
	return wd.Path()
}

func (wd *Directory) Path() string {
	if wd.path == "" {
		d, _ := os.Getwd()
		return d
	}
	return wd.path
}

func (wd *Directory) TouchAll(filenames ...string) ([]string, error) {
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

func (wd *Directory) Touch(filename string) (string, error) {
	fh, err := os.Create(wd.Join(filename))
	if err != nil {
		return filename, err
	}
	return filename, fh.Close()
}

func (wd *Directory) Join(filename string) string {
	return filepath.Join(wd.Path(), filename)
}

func (wd *Directory) RemoveAll() error {
	if wd.Path() == "/" {
		return fmt.Errorf("Cannot remove root directory")
	}
	return os.RemoveAll(wd.Path())
}

// Copy the src file to dest. Both src and dest are considered to be
// inside the working dir.
func (wd *Directory) Copy(dest, src string) (err error) {
	in, err := os.Open(wd.Join(src))
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(wd.Join(dest))
	if err != nil {
		return
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return
}
