#!/bin/sh

if [ "$1" = "docker" ]; then
  cmd="docker build $2 $3 -t longkai/xiaolongtongxue.com:latest ."
  eval $cmd
  exit $?
fi

if [ ! -d "assets/bower_components" ]; then
  # fetch frontend assets if necessary
  cd assets && bower install && cd ..
fi

rev=`git rev-parse --short HEAD`

cmd="go build -ldflags \"-X main.rev=$rev\""

echo "$cmd" && eval $cmd
