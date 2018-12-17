package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gregoryv/workdir"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	wd := workdir.New(path)
	if len(os.Args) > 1 {
		abs, err := filepath.Abs(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		wd = workdir.New(abs)
	}
	wd.LsGit(os.Stdout, true)
}