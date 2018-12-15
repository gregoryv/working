package main

import (
	"github.com/gregoryv/dir"
)

func main() {
	wd := dir.NewPath()
	wd.Ls()
}
