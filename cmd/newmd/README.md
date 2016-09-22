Create New Markdown
===
Auto creating markdown format skelton command line tool. You can find the sample [here][sample].

Download binary [here][dl] or download the `main.go` to build on your own.

## Usage
Run `newmd title`, for example,

```sh
$ newmd "Happy Hacking"
done :)

$ ls -R
happy-hacking/

./happy-hacking:
README.md

$ cat happy-hacking/README.md                       ✱ ◼
Happy Hacking
===
Content goes here...

### EOF
{{yaml fenced code block}} # turncate here... see the sample link
```

[sample]: https://raw.githubusercontent.com/longkai/xiaolongtongxue.com/master/render/testdata/normal.md
[dl]: https://dl.xiaolongtongxue.com/newmd/
