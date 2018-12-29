#!/bin/sh
#
# Jarvis runs this script each time a file changes within the project
#
clear
path=$1
filename=$(basename "$path")
extension="${filename##*.}"

case $extension in
  go)
      gofmt -w $path
      go test github.com/gregoryv/red-rabbit/cursor
      go test github.com/gregoryv/red-rabbit/man
    ;;
  *)
    ;;
esac

