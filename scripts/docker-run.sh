#!/bin/sh

# Always with the lastest docker image
# usage: docker-run.sh [tag] [container-name]

tag=$1
container=$2

if [ -z "$tag" ]; then
  tag="latest"
fi

if [ -z "$container" ]; then
  container="essays"
fi


image="longkai/xiaolongtongxue.com:$tag"

# update docker image if any
echo "pulling $image"
docker pull $image | grep "is up to date" &> /dev/null
uptodate=$?

# stop if it's running
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

# NOTE: don't forget to change your environment!
docker run -d -p 1217:1217 --name=$container -v $HOME/env.yml:/env.yml:ro  -v $HOME/repo:/repo $image

# clean up may goes here...
