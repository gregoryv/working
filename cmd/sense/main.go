package main

import (
	"context"
	"fmt"

	"github.com/gregoryv/working"
)

func main() {
	d := new(working.Directory)
	s := working.NewSensor()
	s.React = func(d *working.Directory, modified ...string) {
		for _, f := range modified {
			fmt.Println(f)
		}
	}
	s.Recursive = true
	d.Watch(context.Background(), s)
}
