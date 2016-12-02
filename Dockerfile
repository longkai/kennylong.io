FROM centos:latest

MAINTAINER longkai <im.longkai@gmail.com>

# Install Git
RUN yum -y install git

# Install Golang
ARG go=1.7.1
ADD https://storage.googleapis.com/golang/go${go}.linux-amd64.tar.gz go.tar.gz
RUN tar -C /usr/local -xzf go.tar.gz && rm -rf go.tar.gz
ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# Install Nodejs&Bower
ARG node=v6.5.0
ADD https://nodejs.org/dist/${node}/node-${node}-linux-x64.tar.xz node.tar.xz
RUN tar -C /usr/local -xf node.tar.xz && mv /usr/local/node* /usr/local/node && rm -rf node.tar.xz
ENV PATH $PATH:/usr/local/node/bin
RUN npm install -g bower && echo '{ "allow_root": true }' > /root/.bowerrc

# This is only a placeholder for check everything is fine. Plz replace yours using build arg at build time or `docker volume` to replace repo dir at runtime
ARG repo=https://github.com/longkai/essays.git
RUN git clone --depth=1 ${repo} /repo

# Compile and Build
ENV workdir $GOPATH/src/github.com/longkai/xiaolongtongxue.com
COPY . $workdir
WORKDIR $workdir
RUN go get ./...
WORKDIR $workdir/assets
RUN bower install
# Replace my fav mdl theme here, you can do it using `docker volume`
ADD https://code.getmdl.io/1.2.1/material.light_green-green.min.css bower_components/material-design-lite/material.min.css
# Predownload Google Fonts
WORKDIR $workdir/cmd/gfdl
RUN go install
WORKDIR $workdir/assets/fonts
RUN gfdl "https://fonts.googleapis.com/css?family=Roboto:regular,bold,italic,thin,light,bolditalic,black,medium&amp;lang=zh" fonts.css
RUN gfdl "https://fonts.googleapis.com/icon?family=Material+Icons" icons.css
# Replace our own which has build ID maybe helpful for CDN cache problem
WORKDIR $workdir
RUN ./build.sh

# Cleanup
RUN npm uninstall -g bower
RUN mv assets templ xiaolongtongxue.com $GOPATH
RUN rm -rf $GOPATH/src $GOPATH/bin $GOPATH/pkg /usr/local/go /usr/local/node*

# Setup
WORKDIR $GOPATH
EXPOSE 1217
VOLUME ["/repo", "/env.yml"]
ENTRYPOINT ["./xiaolongtongxue.com"]
CMD ["/env.yml", "2>&1 | tee /log.txt"] # Let users override mounted configutation file path
