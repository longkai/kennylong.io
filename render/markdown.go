package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/longkai/xiaolongtongxue.com/env"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	ENDPOINT = "https://api.github.com"
)

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
