#!/bin/sh

# usage: docker-run.sh [tag]

tag=$1

if [ -z "$tag" ]; then
  tag="latest"
fi

image="longkai/xiaolongtongxue.com:$tag"

# update docker image if any
echo "pulling $image"
docker pull $image | grep "is up to date" &> /dev/null
uptodate=$?

# stop if it's running
container="sakura"
running=`docker inspect  --format="{{ .State.Running }}" $container` 
if [ $? -eq 0 ]; then
  if [ $running = "true" ]; then
    if [ $uptodate -eq 0 ]; then
      echo "it's already running and up to date, do nothing."
      exit 0
    fi
    echo "stopping $container"
    docker stop $container
  fi

  echo "removing stale $container"
  docker rm $container
fi

docker run -d -p 1217:1217 --name=$container -v $HOME/env.yml:/env.yml:ro $image

# clean up may goes here...
