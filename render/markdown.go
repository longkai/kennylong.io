package render

import (
	"bytes"
	"html/template"
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
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Tags       []string  `json:"tags"`
	Date       time.Time `json:"date"`
	Weather    string    `json:"weather"`
	Summary    string    `json:"summary"`
	Location   string    `json:"location"`
	Background string    `json:"background"`
	Hide       bool      `json:"hide"` //  hide from the list, but still can get will url
	body       []byte    // TODO: avoid this field?
}

// Markdown a rendered *md file.
type Markdown struct {
	Meta
	Older, Newer string
	Body         template.HTML
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
	if m.body, err = linkify(bytes.NewReader(m.body), []byte(m.ID)); err != nil {
		return nil, err
	}
	// don't forget to linkify meta's url
	if len(m.Background) != 0 && !strings.HasPrefix(m.Background, `#`) && !reAbsURL.MatchString(m.Background) {
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

	m.Title = title
	m.body = body
	return parseYAML(bytes.NewReader(_yaml), m)
}

var (
	reAbsURL = regexp.MustCompile(`(?i)((^https?:\/\/)|(^[\/]{1,2}))[^\/]?\S*`)
	reMDLink = regexp.MustCompile(`\[([^\]]*)\]\(([^)"]+)(?: \"([^\"]+)\")?\)`)
)

// linkify makes all the relative links to absolute for easiliy handling
var linkify = func(in io.Reader, prefix []byte) ([]byte, error) {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	ins := reMDLink.FindAllSubmatchIndex(b, -1)
	var i int
	return reMDLink.ReplaceAllFunc(b, func(old []byte) []byte {
		defer func() { i++ }()
		length := ins[i][2*2+1] - ins[i][2*2]     // url len
		offset := ins[i][2*1+1] - ins[i][2*1] + 3 // 3 = [](
		if !reAbsURL.Match(old[offset : offset+length]) {
			// alloc new memory
			buf := make([]byte, 0, len(old)+len(prefix))
			buf = append(buf, old[:offset]...)
			buf = append(buf, prefix...)
			return append(buf, old[offset:]...)
		}
		return old
	}), nil
}
