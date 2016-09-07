xiaolongtongxue.com
===
[![Build Status](https://travis-ci.org/longkai/xiaolongtongxue.com.svg?branch=master)](https://travis-ci.org/longkai/xiaolongtongxue.com)
[![Docker Automated build](https://img.shields.io/docker/automated/jrottenberg/ffmpeg.svg?maxAge=2592000)](https://hub.docker.com/r/longkai/xiaolongtongxue.com/)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![License CC BY 4.0](https://img.shields.io/badge/License-CC%20BY%204.0-lightgrey.svg)](http://creativecommons.org/licenses/by/4.0/)

Frontend and backend source of https://xiaolongtongxue.com

It builds upon **Github Fav Markdown API**, rendering from a plain markdown repo to a nice website. Moreover, it supports **auto update** when you push commits to Github.

It's highly **customizable** and even has a docker image for build-run-ship easily.

## Features
- **Auto** Using Git/Github to keep your writing workflow, when you push your work to Github, your website will sync changes automatically
- **Standard** Github Fav Markdown rendering style and API
- **Fast** Using Non-blocking archetecture, no really a static website but dynamical
- **Docker** Run right from Docker
- **Support CDN** Put all your static stuff to *CDN*(Only tested qiniu)
- **Configurable** You can modify for your needs

## Markdown format Requirement
1. Each doc must have an directory
2. Each doc must ends with `.md`
3. Must have a `EOF` Fenced code block, all the rest has no restrict,

Note the format is(at least one `#`),

```md
### EOF
{{yaml fenced code block}}
```

```yaml
--- sample
background: banner image for this article
date: 2016-01-07T02:50:41+08:00 # must be this format
hide: false # if true this article won't show in the list
location: somewhere 
summary: summary for this article
tags:
  - tag1
  - tag2
  - ...
weather: hey, what's the weather like?
```

## Run with Docker
Run `docker run -d -p 1217:1217 -v /path/to/repo:/repo -v /path/to/env.yaml:/env.yaml:ro longkai/xiaolongtongxue.com` Don't forget to replace your volumes.

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

### Building
1. `go get github.com/longkai/xiaolongtongxue.com && rm $GOPATH/bin/xiaolongtongxue.com`
2. `cd $GOPATH/src/github.com/longkai/xiaolongtongxue.com`
3. `./build.sh`
4. `./xiaolongtongxue.com [/path/to/env.yaml]`

## Configuration
```yaml
--- env.yaml
port: 1217
repo: /repo
hook_secret: Github WebHook secret
access_token: Github Personal access token
meta:
  ga: GA tracker ID
  gf: false # Use Google Fonts, check `templ/include.html`
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
ignores: # NOTE the path is **HTTP Path** format
  - '^/[^/]+\.md$' # ignore *.md in root dir
```

Note if you use docker image with which container has a mounted repo, the `repo` in the `env.yaml` and the docker mount pointer MUST be same.

## CDN Support
I only tested *qiniu CDN* which can fetch then cache your site stuff for a given url. You must set your site url with prefix `/cdn/` to qiniu, then specify the CDN domain in `env.yaml`.

Happy hacking.

[go]: https://golang.org/
[bower]: https://bower.io/
