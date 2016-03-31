xiaolongtongxue.com
===
Frontend and backend source code of https://www.xiaolongtongxue.com 

## Why and What this repo?
[See here][1]

## Architecture
![](seq.png)

## Config
```json
{
  "hook_secret": "your github hook secret",
  "access_token": "your personal github access token",
  "article_repo": "your article git repo",
  "publish_dirs": ["top level", "path", "will", "be", "render", "as", "html"]
}
```

``hook_secret`` is optional, others are required.

## How to build on your own machine?
Download or clone this repo, then [download the binary executable manually][2] and puts them in the same directory.

Or if you have go installed.
```sh
$ git clone https://github.com/longkai/xiaolongtongxue.com.git
$ cd xiaolongtongxue.com
$ go build
```

finally, run it!

```sh
# Usage of ./xiaolongtongxue.com:
#  -conf string
#    	config file path (default "testing_env.json")
#  -port int
#    	http port number (default 1217)
$ ./xiaolongtongxue.com -conf=env.json -port=8080
```

## Note
1. the ``gen`` dir is where the output article html resides, don't modify file under it but you can delete it as you like.
2. the html use google fonts, if you have issue accessing google, you may download the fonts and icons manually in the file ``templ/include.html``, line 29. You may find [this][3] useful.

## License
- code MIT
- content CC BY 4.0

[1]: https://www.xiaolongtongxue.com/memories/2016/go_revamp
[2]: http://dl.xiaolongtongxue.com/xiaolongtongxue/
[3]: https://github.com/longkai/gfdl
