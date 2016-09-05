#!/bin/sh

if [ "$1"="docker" ]; then
  cmd="docker build $2 $3 -t longkai/sakura:latest ."
  eval $cmd
  exit $?
fi

rev=`git rev-parse --short HEAD`

cmd="go build -ldflags \"-X main.rev=$rev\""

echo "$cmd" && eval $cmd
