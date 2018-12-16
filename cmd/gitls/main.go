package main

import (
	"os"

	"github.com/gregoryv/workdir"
)

func main() {
	wd := workdir.New()
	wd.LsGit(os.Stdout)
}
