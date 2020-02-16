package working

import (
	"os"
)

func setup() (wd *Directory, err error) {
	wd, err = TempDir()
	if err != nil {
		return
	}
	os.Chdir(wd.Path())
	wd.MkdirAll("sub/lev", "empty", "ex", "newdir")
	_, err = wd.TouchAll("A", "B", "sub/lev/C", ".hidden", "ex/D")
	if err != nil {
		return
	}

	wd.WriteFile("A", []byte("hello"))
	wd.WriteFile("sub/lev/C", []byte("hello"))
	wd.WriteFile("sub/lev/C", []byte("world"))
	wd.TouchAll("ex/e1", "ex/e2", "newdir/file.txt")
	return
}
