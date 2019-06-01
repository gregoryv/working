#!/bin/bash -e

path=$1
filename=$(basename "$path")
extension="${filename##*.}"

case $extension in
    go)
	goimports -w $path
        gofmt -w $path
        ;;
esac
go test -coverprofile /tmp/c.out .
uncover /tmp/c.out
