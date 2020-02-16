package working

import (
	"bytes"
	"os"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestChmod(t *testing.T) {
	d := new(Directory)
	d.Temporary()
	f := "file"
	d.Touch(f)
	err := d.Chmod(f, 0400) // read only
	if err != nil {
		t.Error(d, err)
	}
	d.RemoveAll()
}

func TestLoad(t *testing.T) {
	d := new(Directory)
	d.Temporary()
	_, err := d.Load("nosuchfile")
	if err == nil {
		t.Error("should fail when loading nonexisting file")
	}

	exp := "hello"
	d.WriteFile("x", []byte(exp))
	body, err := d.Load("x")
	if err != nil {
		t.Error(err)
	}
	got := string(body)
	assert := asserter.New(t)
	assert().Equals(got, exp)
	d.RemoveAll()
}

func TestIsEmpty(t *testing.T) {
	d := new(Directory)
	d.Temporary()
	if !d.IsEmpty("") {
		t.Error("Expected new temporary directory to be empty")
	}
	d.TouchAll("k")
	if d.IsEmpty("") {
		t.Error("Dir with contents should not be empty")
	}

	d.RemoveAll()
	if d.IsEmpty("") {
		t.Error("IsEmpty should be false for non existing")
	}

}

func TestRemoveAll(t *testing.T) {
	d := new(Directory)
	d.SetPath("/")
	err := d.RemoveAll() // :-)
	if err == nil {
		t.Fatal("Well we've probably erased the entire disk")
	}
}

func TestTempDir_error(t *testing.T) {
	os.Setenv("TMPDIR", "/_no_such_dir")
	defer os.Setenv("TMPDIR", "/tmp")
	d := new(Directory)
	err := d.Temporary()
	if err == nil {
		t.Fail()
	}
}

func Test_Ls_error(t *testing.T) {
	d := new(Directory)
	d.Temporary()
	d.RemoveAll()
	err := d.Ls(nil)
	if err == nil {
		t.Fail()
	}
}

func Test_Ls(t *testing.T) {
	d := new(Directory)
	d.Temporary()
	setup(d)
	defer d.RemoveAll()

	w := bytes.NewBufferString("")
	d.Ls(w)
	got := w.String()
	exp := `A
B
empty/
ex/
newdir/
sub/
`
	if exp != got {
		t.Errorf("Expected \n'%s'\ngot \n'%s'", exp, got)
	}
}

func setup(d *Directory) {
	os.Chdir(d.Path())
	d.MkdirAll("sub/lev", "empty", "ex", "newdir")
	_, err := d.TouchAll("A", "B", "sub/lev/C", ".hidden", "ex/D")
	if err != nil {
		panic(err)
	}
	d.WriteFile("A", []byte("hello"))
	d.WriteFile("sub/lev/C", []byte("hello"))
	d.WriteFile("sub/lev/C", []byte("world"))
	d.TouchAll("ex/e1", "ex/e2", "newdir/file.txt")
}

func Test_String(t *testing.T) {
	got := new(Directory).Path()
	exp, _ := os.Getwd()
	if exp != got {
		t.Errorf("Expected %q, got %q", exp, got)
	}
}

func TestTouch(t *testing.T) {
	d := new(Directory)
	d.Temporary()
	defer d.RemoveAll()
	d.MkdirAll("x")
	os.Chmod(d.Join("x"), 0000)
	_, err := d.TouchAll("x")
	if err == nil {
		t.Error("Expected to fail")
	}
	// Cleanup
	os.Chmod(d.Join("x"), 0644)
}

func TestMkdirAll(t *testing.T) {
	d := new(Directory)
	d.Temporary()
	defer d.RemoveAll()
	os.Chmod(d.Path(), 0000)
	err := d.MkdirAll("hepp")
	if err == nil {
		t.Error("Expected to fail")
	}
	os.Chmod(d.Join("x"), 0644)
}

func TestCopy_errors(t *testing.T) {
	d := new(Directory)
	d.Temporary()
	defer d.RemoveAll()
	err := d.Copy("dest", "nosuchfile")
	if err == nil {
		t.Error("Should fail to copy when src doesn't exist")
	}
	a := "a.file"
	b := "b.file"
	body := []byte("content")
	d.WriteFile(a, body)
	d.WriteFile(b, body)
	err = os.Chmod(d.Join(a), 0000)
	if err != nil {
		t.Error(err)
	}
	err = d.Copy(a, b)
	if err == nil {
		t.Error("Should fail if cannot write destination")
	}
	os.Chmod(d.Join(a), 0644)
}
