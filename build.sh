#!/bin/sh

# TODO: since we've imposed docker build automate, this should be deleted.
if [ "$1" = "docker" ]; then
  cmd="docker build $2 $3 -t longkai/xiaolongtongxue.com:`git rev-parse --abbrev-ref HEAD` ."
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

cmd="go build -ldflags \"\
  -X github.com/longkai/xiaolongtongxue.com/config.v=`git rev-parse --short HEAD` \
  -X github.com/longkai/xiaolongtongxue.com/config.b=`git rev-parse --abbrev-ref HEAD` \
  \""

echo "$cmd" && eval $cmd
