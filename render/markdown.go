package render

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/longkai/xiaolongtongxue.com/config"
)

// Meta metadata for the markdown.
type Meta struct {
	ID           string
	body         []byte    // TODO: avoid this field?
	Title        string    `json:"title"`
	Tags         []string  `json:"tags"`
	Date         time.Time `json:"date"`
	Weather      string    `json:"weather"`
	Summary      string    `json:"summary"`
	Location     string    `json:"location"`
	Background   string    `json:"background"`
	RenderOption int       `json:"render_option"`
}

// Markdown a rendered *md file.
type Markdown struct {
	Meta
	Prev, Next string
	Body       template.HTML
}

var parseJSON = func(in io.Reader, v interface{}) error {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}
	return json.Unmarshal(b, v)
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
	title, body, _json, err := parse(in)
	if err != nil {
		return nil, err
	}

	m := new(Meta)
	m.Title = title
	m.body = body
	if err = parseJSON(bytes.NewReader(_json), m); err != nil {
		log.Print(err)
	}
	return m, nil
}
