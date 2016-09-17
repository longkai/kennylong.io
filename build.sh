#!/bin/sh

if [ "$1" = "docker" ]; then
  cmd="docker build $2 $3 -t longkai/xiaolongtongxue.com:latest ."
  eval $cmd
  exit $?
fi

# fetch frontend assets if necessary
if [ ! -d "assets/bower_components" ]; then
  cd assets && bower install && cd ..
fi

# test existence of Google Fonts 
if [ ! -d "assets/fonts" ]; then
  echo "Google fonts have not yet been downloaded locally. Checkout 'templ/include.html' or 'cmd/gfdl' for more information."
fi

rev=`git rev-parse --short HEAD`

cmd="go build -ldflags \"-X main.rev=$rev\""

echo "$cmd" && eval $cmd
