package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gregoryv/dir"
)

func main() {
	wd := dir.NewPath()
	wd.Ls()

	status, err := exec.Command("git", "status", "-s", "-C", wd.Root).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s", string(status))
}
