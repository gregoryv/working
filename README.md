[![Build Status](https://travis-ci.org/gregoryv/working.svg?branch=master)](https://travis-ci.org/gregoryv/working)
[![codecov](https://codecov.io/gh/gregoryv/working/branch/master/graph/badge.svg)](https://codecov.io/gh/gregoryv/working)


[working](https://godoc.org/github.com/gregoryv/working) - Package defines Directory type for easy file operations

This package is most useful if you need to do multiple file operations
within one directory. Eg. something like

    cd /tmp/dir
	mkdir child1 child2
	touch child1/A child2/B
	# do some stuff
	cd
	rm -rf /tmp/dir

Using this package you would

    d := new(working.Directory)
	d.SetPath("/tmp/dir")
	d.MkdirAll("child1", "child2")
	d.TouchAll("child1/A", "child2/B")
	// do stuff
	d.RemoveAll()
