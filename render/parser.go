package render

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
)

var (
	re = regexp.MustCompile("\\s*(?:#+\\s*)?(?P<title>\\S[\\S ]*\\S)?(?:\\s*[-=]+\\s*)?(?P<body>[\\S\\s]+?)(#+\\s*EOF\\s+```json\\s*(?P<json>[\\S\\s]+)```)(?P<links>[\\S\\s]*)")
)

// parse parsing the markdown then extracting the metas. If no title matches, the first non-blank line will be used as title. **Note title and JSON(required) will be stripped from the body**.
func parse(in io.Reader) (title string, body, json []byte, err error) {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return
	}

	indices := re.FindSubmatchIndex(b)
	if len(indices)/2 != 6 {
		err = fmt.Errorf("parse regexp not matched")
		return
	}

	title = string(b[indices[2*1]:indices[2*1+1]]) // <title>

	// avoid allocating memory, do some tricks here
	lo1, hi1, lo2, hi2 := indices[2*3], indices[2*3+1], indices[2*5], indices[2*5+1]
	tail := b[lo1:hi2]
	reverse(tail[:hi1-lo1])
	reverse(tail[lo2-lo1:])
	reverse(tail)

	json = tail[hi2-hi1+lo2-hi1+indices[2*4]-lo1 : len(tail)-hi1+indices[2*4+1]] // <json> strip from wrapper
	body = b[indices[2*2] : indices[2*2+1]+hi2-lo2]                              // <body> + <links>
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
