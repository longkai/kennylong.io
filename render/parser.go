package render

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
)

var (
	re = regexp.MustCompile(`(?s)\s*(?:#+\s*)?(?P<title>\w[\w ]*\w)?(?:\s*[-=]+\s*)?(?P<body>[^#$]*)(#+\s*EOF\s*` + "```json" + `\s*(?P<json>.*)` + "```)?(?P<links>.*)")
)

// parse parsing the markdown then extracting the metas. If no title matches, the first non-blank line will be used as title. **Note title and JSON(if any) will be stripped from the body**.
func parse(in io.Reader) (title string, body, json []byte, err error) {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return
	}
	indices := re.FindSubmatchIndex(b)
	if len(indices)/2 == 6 {
		title = string(b[indices[2*1]:indices[2*1+1]]) // <title>

		if indices[2*3] != -1 { // has meta JSON block
			// avoid allocating memory, do some tricks here
			lo1, hi1, lo2, hi2 := indices[2*3], indices[2*3+1], indices[2*5], indices[2*5+1]
			tail := b[lo1:hi2]
			reverse(tail[:hi1-lo1])
			reverse(tail[lo2-lo1:])
			reverse(tail)

			json = tail[hi2-hi1+lo2-hi1+indices[2*4]-lo1 : len(tail)-hi1+indices[2*4+1]] // <json> strip from wrapper
			body = b[indices[2*2] : indices[2*2+1]+hi2-lo2]                              // <body> + <links>
		} else {
			// has no meta JSON block, will match all the rest
			body = b[indices[2*2]:]
		}
		return
	}
	err = fmt.Errorf("not matched")
	return
}

func reverse(b []byte) {
	lo, hi := 0, len(b)-1
	for lo < hi {
		b[lo], b[hi] = b[hi], b[lo]
		lo++
		hi--
	}
}
