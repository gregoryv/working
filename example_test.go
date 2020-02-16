package working

func ExampleTempDir() {
	d := new(Directory)
	d.Temporary()
	d.WriteFile("hello.txt", []byte("Hello, working directory!"))
	d.TouchAll("a.md", "b.sh")
	d.Ls(nil)
	d.RemoveAll()
	// output:
	// a.md
	// b.sh
	// hello.txt

}

func ExampleWorkDir_Copy() {
	d := new(Directory)
	d.Temporary()
	defer d.RemoveAll()
	d.WriteFile("src", []byte("hello"))
	d.Copy("dest", "src")
	d.Ls(nil)
	// output:
	// dest
	// src
}
