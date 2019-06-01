package workdir

func ExampleTempDir() {
	wd, _ := TempDir()
	wd.WriteFile("hello.txt", []byte("Hello, working directory!"))
	wd.TouchAll("a.md", "b.sh")
	wd.Ls(nil)
	wd.RemoveAll()
	// output:
	// a.md
	// b.sh
	// hello.txt

}

func ExampleWorkDir_Copy() {
	wd, _ := TempDir()
	defer wd.RemoveAll()
	wd.WriteFile("src", []byte("hello"))
	wd.Copy("dest", "src")
	wd.Ls(nil)
	// output:
	// dest
	// src
}
