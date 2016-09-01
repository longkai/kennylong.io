package render

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/longkai/xiaolongtongxue.com/config"
)

// Meta metadata for the markdown.
type Meta struct {
	ID         string
	body       []byte    // TODO: avoid this field?
	Title      string    `yaml:"title"`
	Tags       []string  `yaml:"tags"`
	Date       time.Time `yaml:"date"`
	Weather    string    `yaml:"weather"`
	Summary    string    `yaml:"summary"`
	Location   string    `yaml:"location"`
	Background string    `yaml:"background"`
	Hide       bool      `yaml:"hide"`
}

// Markdown a rendered *md file.
type Markdown struct {
	Meta
	Prev, Next string
	Body       template.HTML
}

var parseYAML = func(in io.Reader, v interface{}) error {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, v)
}

func parseID(path string) string {
	base := config.Env.ArticleRepo
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

func parseMd(in io.Reader) (*Meta, error) {
	title, body, _yaml, err := parse(in)
	if err != nil {
		return nil, err
	}

	m := new(Meta)
	m.Title = title
	m.body = body
	err = parseYAML(bytes.NewReader(_yaml), m)
	return m, err
}
