package render

import (
	"bytes"
	"encoding/json"
	"github.com/longkai/xiaolongtongxue.com/env"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	def  = iota // render and appear in the article list
	keep        // render but disappear in the article list
	skip        // no render

	fenced_block = "```"
)

var (
	titleRegexp = regexp.MustCompile(`([\S ]+)\s*=+\s+`)
	metaRegexp  = regexp.MustCompile(`#+\s*EOF\s+` + fenced_block + `json\s+([\s\S]*)\s+` + fenced_block)
)

type markdownMeta struct {
	Id           string    `json:"id"`
	Title        string    `json:"title"`
	Tags         []string  `json:"tags"`
	Date         time.Time `json:"date"`
	Weather      string    `json:"weather"`
	Summary      string    `json:"summary"`
	Location     string    `json:"location"`
	Background   string    `json:"background"`
	RenderOption int       `json:"render_option"`
}

type markdown struct {
	markdownMeta
	Text string
	Html template.HTML
}

func newMarkdown(src string) (*markdown, error) {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}

	title, b := parseTitle(b) // separate title block
	text, meta := separateTextAndMeta(b)
	meta.Id = trimLastSegment(src[len(env.Config().ArticleRepo):])
	if meta.Title == "" { // if title not provided in json meta, use the markdown body if has
		meta.Title = title
	}
	m := new(markdown)
	m.Text = text
	m.markdownMeta = meta
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

func separateTextAndMeta(slice []byte) (string, markdownMeta) {
	meta := markdownMeta{}
	result := metaRegexp.FindSubmatch(slice)
	if result == nil {
		return string(slice), meta
	}
	err := json.Unmarshal(result[1], &meta)
	if err != nil {
		log.Printf("separateTextAndMeta %s fail, %v\n", result[1], err)
		return "", meta
	}
	// drop the json code block
	slice = bytes.Replace(slice, result[0], []byte(""), -1)
	return string(slice), meta
}

// trim last segment of the file path
// i.e. a/b/c/d.ext to a/b/c
func trimLastSegment(s string) string {
	i := strings.LastIndexFunc(s, func(r rune) bool {
		return r == filepath.Separator
	})

	if i != -1 {
		return s[:i]
	}
	return s
}

// helper functions for template usage
func DaysAgo(t time.Time) int { return int(time.Since(t).Hours() / 24) }

func Format(t time.Time) string { return t.Format(time.RFC1123Z) }

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
