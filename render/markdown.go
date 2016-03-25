package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/longkai/xiaolongtongxue.com/env"
	"html/template"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

const (
	fenced_block = "```"
)

var (
	titleRegexp = regexp.MustCompile(`([\S ]+)\s*=+\s+`)
	metaRegexp  = regexp.MustCompile(`#+\s*EOF\s+` + fenced_block + `json\s+([\s\S]*)\s+` + fenced_block)
)

type MarkdownMeta struct {
	Id         string    `json:"id"`
	Title      string    `json:"title"`
	Tags       []string  `json:"tags"`
	Reserved   bool      `json:"reserved"`
	Date       time.Time `json:"date"`
	Weather    string    `json:"weather"`
	Summary    string    `json:"summary"`
	Location   string    `json:"location"`
	Background string    `json:"background"`
}

type Markdown struct {
	Text string
	Html template.HTML
	MarkdownMeta
}

func NewMarkdown(src string) (*Markdown, error) {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}

	title, b := parseTitle(b) // seperate title block
	text, meta := separateTextAndMeta(b)
	meta.Id = trimBasename(src[len(env.Config().ArticleRepo):])
	if meta.Title == "" { // if title not provided in json meta
		meta.Title = title
	}
	m := new(Markdown)
	m.Text = text
	m.MarkdownMeta = meta
	return m, nil
}

func parseTitle(slice []byte) (string, []byte) {
	result := titleRegexp.FindSubmatch(slice)
	if result == nil {
		return "", slice
	}
	t := string(bytes.TrimSpace(result[1]))
	slice = slice[len(result[0]):]
	return t, slice
}

func separateTextAndMeta(slice []byte) (string, MarkdownMeta) {
	meta := MarkdownMeta{}
	result := metaRegexp.FindSubmatch(slice)
	if result == nil {
		return string(slice), meta
	}
	err := json.Unmarshal(result[1], &meta)
	if err != nil {
		fmt.Println(err)
		return "", meta
	}
	// drop the json code block
	slice = bytes.Replace(slice, result[0], []byte(""), -1)
	return string(slice), meta
}

// trim the basename of the file
func trimBasename(s string) string {
	j := -1
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			j = i
			break
		}
	}
	if j != -1 {
		return s[:j]
	}
	return s
}

// helper functions for template usage
func DaysAgo(t time.Time) int { return int(time.Since(t).Hours() / 24) }

func HasColor(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] == '#'
}

func HasImage(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] != '#'
}

func Tags(m []string) string {
	s := ""
	for i, v := range m {
		s += "#" + v
		if i != len(m)-1 {
			s += ", "
		}
	}
	return s
}

func IsRelImage(s string) bool {
	if len(s) == 0 {
		return false
	}
	if s[0] == '/' || strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return false
	}
	return true
}
