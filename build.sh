#!/bin/sh

rev=`git rev-parse --short HEAD`

cmd="go build -ldflags \"-X main.rev=$rev\""

echo "$cmd" && eval $cmd
