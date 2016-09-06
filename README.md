xiaolongtongxue.com
===
[![Build Status](https://travis-ci.org/longkai/xiaolongtongxue.com.svg?branch=master)](https://travis-ci.org/longkai/xiaolongtongxue.com)
[![Docker Automated buil](https://img.shields.io/docker/automated/jrottenberg/ffmpeg.svg?maxAge=2592000)](https://hub.docker.com/r/longkai/xiaolongtongxue.com/)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![License CC BY 4.0](https://img.shields.io/badge/License-CC%20BY%204.0-lightgrey.svg)](http://creativecommons.org/licenses/by/4.0/)

Frontend and backend source code of https://xiaolongtongxue.com

It builds upon **Github Fav Markdown API**, rendering from a plain markdown repo to a nice website. Moreover, it supports **auto update** when you push commmits to Github.

It's highly **customizable** and even has a docker image for build-run-ship easily.

## Run with Docker
`docker run -d -p 1217 -v /path/to/repo:/repo -v /path/to/env.yaml:/env.yaml longkai/xiaolongtongxue.com`

Don't forget mount your volumes to the container.

Or, if you prefer `docker-compose`, modify for your needs,

```yaml
sakura:
  image: longkai/xiaolongtongxue.com
  ports:
    - "1217:1217"
  volumes:
    - /path/to/env.yaml:/env.yaml:ro
    - /path/to/repo:/repo
```

then run `docker-compose up -d`

## Build Manually
### Pre-requisite
- [golang][go] >= 1.7
- [bower][bower]

1. `git clone https://github.com/longkai/xiaolongtongxue.com.git`
2. `./build.sh`
3. modify configurations, see below
4. `./xiaolongtongxue.com [/path/to/env.yaml]`

## Configuration
```yaml
port: 1217
repo: /repo
hook_secret: Github WebHook secret
access_token: Github Personal access token
meta:
  ga: GA tracker ID
  #cdn: CDN domain # currently only tested qiniu
  domain: domain.com # required only if using CDN
  bio: something about you
  link: other link about you
  lang: zh
  name: your name
  title: page title
  mail: you@somewhere
  domain: domain.com # optional, for multiple sub-domain tracking
  github: your Github link if nay
  medium: medium repo if any
  twitter: twitter link if any
  instagram: ins link if any
  stackoverflow: stackoverflow link if any
ignores: # note the path is HTTP Path format
  - '^/[^/]+\.md$' # ignore *.md in root dir
```

Note if you use docker image which mount your repo to the container, the `repo` in the `env.yaml` and the docker mount pointer MUST be same.

Happy hacking.
