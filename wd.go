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

func New(path string) WorkDir {
	return WorkDir(path)
}

type WorkDir string

func (wd WorkDir) Chmod(filename string, mode os.FileMode) error {
	return os.Chmod(wd.Join(filename), mode)
}

// IsEmpty returns true if the dir is empty, false otherwise.
// Use empty string to check the workdir itself.
func (wd WorkDir) IsEmpty(dir string) bool {
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
func TempDir() (WorkDir, error) {
	tmpPath, err := ioutil.TempDir("", "workdir")
	if err != nil {
		return WorkDir(""), err
	}
	return WorkDir(tmpPath), nil
}

// WriteFile creates/writes over the file with mode 0644
func (wd WorkDir) WriteFile(file string, data []byte) error {
	return ioutil.WriteFile(wd.Join(file), data, 0644)
}

// ReadAll loads the given file like ioutil.ReadAll
func (wd WorkDir) Load(file string) ([]byte, error) {
	fh, err := os.Open(wd.Join(file))
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	return ioutil.ReadAll(fh)
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
	if string(wd) == "/" {
		return fmt.Errorf("Cannot remove root directory")
	}
	return os.RemoveAll(wd.String())
}

// Copy the src file to dest. Both src and dest are considered to be
// inside the working dir.
func (wd WorkDir) Copy(dest, src string) (err error) {
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
