Google Fonts Downloader
===
Have your site used Google Fonts? If so, for some reasons(e.g. firewall, etc.) you may want to pre-download all the assets on your site.

This tool helps you automate it.

## Install
Download `main.go` then run `go build` or download directly [here][dl].

## Usage
```sh
$ gfdl src dst
# gfdl "https://fonts.googleapis.com/css?family=Roboto:regular,bold,italic,thin,light,bolditalic,black,medium&amp;lang=zh" path/to/fonts.css
```

Note all the fonts assets will place on the same folder as the css file.

## License
MIT

[dl]: https://dl.xiaolongtongxue.com/gfdl/
