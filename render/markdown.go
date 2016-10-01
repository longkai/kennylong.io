package render

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/longkai/xiaolongtongxue.com/config"
)

// Meta metadata for the markdown.
type Meta struct {
	ID         string      `json:"id,omitempty" yaml:"-"`
	Title      string      `json:"title" yaml:",omitempty"`
	Tags       []string    `json:"tags"`
	Date       time.Time   `json:"date"`
	Weather    string      `json:"weather"`
	Summary    string      `json:"summary"`
	Location   string      `json:"location"`
	Background string      `json:"background"`
	License    string      `json:"license"`    // default all-rights-reserved
	Hide       bool        `json:"hide"`       // hide from the list, but still can get will url
	Body       interface{} `json:"-" yaml:"-"` // initilized as []byte, then render it as template.HTML
}

// Markdown a rendered *md file.
type Markdown struct {
	Meta
	Older, Newer string
}

var parseYAML = func(in io.Reader, v interface{}) error {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, v)
}

// parseID ensures start and end with `/` just same as HTTP RequestURI
var parseID = func(path string) string {
	base := config.Env.Repo
	if !strings.HasPrefix(path, base) {
		return ""
	}
	var id string
	if strings.HasSuffix(base, "/") {
		id = path[strings.LastIndexByte(base, '/'):]
	} else {
		id = path[len(base):]
	}
	// trim last segment
	return id[:strings.LastIndexByte(id, '/')+1]
}

func parseMD(path string) (*Meta, error) {
	return ParseMD(path, "")
}

// ParseMD a markdown doc for a given path, baseURL is used for linkify all the relative links in the doc to absoluate for some other reasons(e.g., let 3rd sync its resources). If `basePrefix` is empty, the doc's ID will be used.
func ParseMD(path, origin string) (*Meta, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m := new(Meta)
	if err = unmarshal(f, m); err != nil {
		return nil, err
	}
	m.ID = parseID(path)
	if origin != "" {
		origin += m.ID
		if m.Body, err = linkifyMD(bytes.NewReader(m.Body.([]byte)), []byte(origin)); err != nil {
			return nil, err
		}
	}
	// don't forget to linkify meta's url
	if len(m.Background) != 0 && m.Background[0] != '#' && !reAbsURL.MatchString(m.Background) {
		// not color bg
		m.Background = m.ID + m.Background
	}
	return m, nil
}

func unmarshal(in io.Reader, m *Meta) error {
	title, body, _yaml, err := parse(in)
	if err != nil {
		return err
	}

	m.Title = title                   // defaults with parsed title if not specified in YAML
	m.License = "all-rights-reserved" // default
	m.Body = body
	return parseYAML(bytes.NewReader(_yaml), m)
}

var (
	reAbsURL   = regexp.MustCompile(`(?i)((^https?:\/\/)|(^[\/]{1,2}))[^\/]?\S*`)
	reMDLink   = regexp.MustCompile(`\[([^\]]*)\]\(([^)"]+)(?: \"([^\"]+)\")?\)`)
	reHTMLLink = regexp.MustCompile(`(?i)<(a|img)[\s\S]+?(href|src)=['"](\S*\.\S+)['"][\s\S]*?>`)
)

// IsAbsURL test url is absolute or not.
func IsAbsURL(url string) bool { return reAbsURL.MatchString(url) }

type indexer func(indices []int) (offset, length int)

func linkify(in io.Reader, prefix []byte, re *regexp.Regexp, f indexer) ([]byte, error) {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	// normalize the prefix
	if prefix[len(prefix)-1] != '/' {
		prefix = append(prefix, '/')
	}

	indices := re.FindAllSubmatchIndex(b, -1)
	var i int
	return re.ReplaceAllFunc(b, func(old []byte) []byte {
		defer func() { i++ }()
		offset, length := f(indices[i])
		if !reAbsURL.Match(old[offset : offset+length]) {
			// alloc new memory
			buf := make([]byte, len(old)+len(prefix))
			copy(buf, old[:offset])
			copy(buf[offset:], prefix)
			copy(buf[offset+len(prefix):], old[offset:])
			return buf
		}
		return old
	}), nil
}

// linkify makes all the relative links in markdown to absolute for easy handling
var linkifyMD = func(in io.Reader, prefix []byte) ([]byte, error) {
	f := func(indices []int) (offset, length int) {
		offset = indices[2*2] - indices[0]     // group 2 is the link
		length = indices[2*2+1] - indices[2*2] // link len
		return
	}
	return linkify(in, prefix, reMDLink, f)
}

// linkify makes all the relative resources links(must ends with .xxx) in HTML to absolute for easy handling
var linkifyHTML = func(in io.Reader, prefix []byte) ([]byte, error) {
	f := func(indices []int) (offset, length int) {
		offset = indices[2*3] - indices[0]     // group 3 is the link
		length = indices[2*3+1] - indices[2*3] // link len
		return
	}
	return linkify(in, prefix, reHTMLLink, f)
}
