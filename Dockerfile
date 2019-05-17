FROM alpine:latest

MAINTAINER longkai <im.longkai@gmail.com>

ENV GOPATH /tmp
ENV PATH $PATH:$GOPATH/bin
ENV SRC /tmp/src/github.com/longkai/xiaolongtongxue.com

ARG branch=master

RUN apk add --no-cache git && \
  apk add --no-cache --virtual .build-deps \
                                  build-base \
                                  linux-headers \
                                  curl \
                                  go \
                                  nodejs-npm \
                                  nodejs && \
  runDeps="$( \
    mkdir -p $SRC && cd $SRC/.. && \
    git clone --branch $branch --depth=1 https://github.com/longkai/xiaolongtongxue.com.git && \
    cd $SRC/assets/ && \
    npm install && \
    rm -rf /root/.[a-zA-Z]* && \
    curl https://code.getmdl.io/1.3.0/material.light_green-green.min.css > node_modules/material-design-lite/material.min.css && \
    cd $SRC/cmd/gfdl/ && \
    go install && \
    cd $SRC/assets/ && \
    gfdl 'https://fonts.googleapis.com/css?family=Roboto:regular,bold,italic,thin,light,bolditalic,black,medium&amp;lang=zh' fonts/fonts.css && \
    gfdl 'https://fonts.googleapis.com/icon?family=Material+Icons' fonts/icons.css && \
    cd $SRC && \
    go get ./... && \
    ./build.sh && \
    mv xiaolongtongxue.com templ/ assets / \
  )" && \
  apk add --no-cache --virtual .run-deps $runDeps && \
  apk del .build-deps && \
  rm -rf /tmp/*


EXPOSE 1217
VOLUME ["/repo", "/conf.yml", "/log.txt"]
CMD /xiaolongtongxue.com /conf.yml 2>&1 | tee /log.txt
