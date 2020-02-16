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

func (d *Directory) SetPath(path string) { d.path = path }

func (d *Directory) Chmod(filename string, mode os.FileMode) error {
	return os.Chmod(d.Join(filename), mode)
}

// IsEmpty returns true if the dir is empty, false otherwise.
// Use empty string to check the workdir itself.
func (d *Directory) IsEmpty(dir string) bool {
	name := d.Join(dir)
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
func (d *Directory) Ls(w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	path := d.Path()
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
func (d *Directory) Temporary() error {
	tmp, err := ioutil.TempDir("", "workdir")
	if err != nil {
		return err
	}
	d.SetPath(tmp)
	return nil
}

// WriteFile creates/writes over the file with mode 0644
func (d *Directory) WriteFile(file string, data []byte) error {
	return ioutil.WriteFile(d.Join(file), data, 0644)
}

// ReadAll loads the given file like ioutil.ReadAll
func (d *Directory) Load(file string) ([]byte, error) {
	fh, err := os.Open(d.Join(file))
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	return ioutil.ReadAll(fh)
}

func (d *Directory) MkdirAll(subDirs ...string) error {
	for _, sub := range subDirs {
		err := os.MkdirAll(filepath.Join(d.Path(), sub), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Directory) String() string {
	return d.Path()
}

func (d *Directory) Path() string {
	if d.path == "" {
		d, _ := os.Getwd()
		return d
	}
	return d.path
}

func (d *Directory) TouchAll(filenames ...string) ([]string, error) {
	files := make([]string, len(filenames))
	for i, name := range filenames {
		name, err := d.Touch(name)
		if err != nil {
			return files, err
		}
		files[i] = name
	}
	return files, nil
}

func (d *Directory) Touch(filename string) (string, error) {
	fh, err := os.Create(d.Join(filename))
	if err != nil {
		return filename, err
	}
	return filename, fh.Close()
}

func (d *Directory) Join(filename string) string {
	return filepath.Join(d.Path(), filename)
}

func (d *Directory) RemoveAll() error {
	if d.Path() == "/" {
		return fmt.Errorf("Cannot remove root directory")
	}
	return os.RemoveAll(d.Path())
}

// Copy the src file to dest. Both src and dest are considered to be
// inside the working dir.
func (d *Directory) Copy(dest, src string) (err error) {
	in, err := os.Open(d.Join(src))
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(d.Join(dest))
	if err != nil {
		return
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return
}
