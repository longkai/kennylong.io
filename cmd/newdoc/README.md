# Create New Markup Documents for Your Site

A command line tool for auto generating your favorite skeleton format. You can find the sample [here][sample].

Currently supports *markdown* and *org-mode*.

## Install

```sh
$ go get github.com/longkai/xiaolongtongxue.com/tree/master/cmd/newdoc
```

## Usage

Run `newdoc title [md | org]`, for example,

```sh
$ newdoc "Happy Hacking"
done :)

$ ls -R
happy-hacking/

./happy-hacking:
README.md

$ cat happy-hacking/README.md
Happy Hacking
===
Content goes here...

## EOF
{{yaml fenced code block}} # truncate here... See the sample link
```

Default type is org-mode.

[sample]: https://raw.githubusercontent.com/longkai/xiaolongtongxue.com/master/render/testdata/normal.md
