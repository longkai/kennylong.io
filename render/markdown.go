package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/longkai/xiaolongtongxue.com/env"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

const (
	ENDPOINT   = "https://api.github.com"
	CODE_BLOCK = "```"
)

var (
	titleRegexp = regexp.MustCompile(`([\S ]+)\s*=+\s+`)
	metaRegexp  = regexp.MustCompile(`#+\s*EOF\s+` + CODE_BLOCK + `json\s+([\s\S]*)\s+` + CODE_BLOCK)
)

type MarkdownMeta struct {
	Tags     []string  `json:"tags"`
	Location string    `json:"location"`
	Weather  string    `json:"weather"`
	Publish  bool      `json:"publish"`
	Date     time.Time `json:"date"`
}

type Markdown struct {
	Text  string
	Title string
	MarkdownMeta
}

func NewMarkdown(src string) (*Markdown, error) {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}
	// the first line is the title
	title, b := parseTitle(b)

	text, meta := separateTextAndMeta(b)
	m := new(Markdown)
	m.Title = title
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

type MarkdownRender struct {
	Text    string `json:"text"`
	Mode    string `json:"mode"`
	Context string `json:"context"`

	reader io.Reader
}

func New(text string) *MarkdownRender {
	return &MarkdownRender{text, "gfm", "github/longkai", nil}
}

func (m *MarkdownRender) Read(p []byte) (n int, err error) {
	if m.reader == nil {
		bs, err := json.Marshal(m)
		if err != nil {
			return 0, err
		}
		m.reader = bytes.NewReader(bs)
	}
	return m.reader.Read(p)
}

// Render the makrdown to html []byte with Github API.
func (m *MarkdownRender) Render() ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, ENDPOINT+"/markdown", m)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("token %s", env.Config().AccessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// the bytes may http error content
	result, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("http StatusCode %d", resp.StatusCode)
	}
	return result, err
}
