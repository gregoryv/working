# backstejg

go library for interacting with [stejg](http://stejg.7de.se)

## Getting started

    go install github.com/gregoryv/backstejg/text/texter
    texter file.md

	go install github.com/gregoryv/backstejg/spinner
	spinner -x 200 -y 200 -r 40 -fc yellow
	
## Write your own actors

    go get github.com/gregoryv/backstejg
   
    import "github.com/gregoryv/backstejg/act"
   

Check out eg. the texter source code for how to send events to stejg.
