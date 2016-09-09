package medium

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/longkai/xiaolongtongxue.com/render"
)

var (
	uid      string
	token    string
	endpoint string

	ready = make(chan struct{})
)

var parse = func(path string) (*render.Meta, error) {
	return render.ParseMD(path, endpoint)
}

var provideURL = func(path string) string {
	return "https://api.medium.com/v1" + path
}

// Init meidum post service, it must be called before any othter functions in this package.
func Init(_token, _endpoint string) {
	if _token == "" || _endpoint == "" {
		log.Fatalf("empty meidum token %q or endpoint %q, aborting", _token, _endpoint)
	}
	token, endpoint = _token, strings.TrimRight(_endpoint, "/") // trim right for pretty URL
	go func() {
		try := 0
		for try < 3 { // retry max 3 times
			try++
			if id, err := me(); err != nil {
				log.Printf("me() fail: %v, tried %d times", err, try)
			} else {
				uid = id
				break
			}
		}
		close(ready)
	}()
}

// Post from the given path
func Post(path string) error {
	<-ready
	if uid == "" {
		return fmt.Errorf("Post(%q) without uid, aborting", path)
	}

	m, err := parse(path)
	if err != nil {
		return err
	}
	b, ok := m.Body.([]byte)
	if !ok {
		return fmt.Errorf("never happen, the parsed markdown is not `[]byte`")
	}
	// always `markdown` and `public` since you have publish to your site
	p := &payload{Title: m.Title, Tags: m.Tags, License: m.License, Format: "markdown", Status: "public"}
	p.CanonicalURL = endpoint + "/" + strings.TrimLeft(m.ID, "/")
	// title has been stripped by parser, however, medium needs it, so we have to prepend it. see their doc at: https://github.com/Medium/medium-api-docs#33-posts
	p.Content = fmt.Sprintf("%s\n===\n%s", m.Title, b)

	if _, err := p.post(); err != nil {
		return err
	}
	return nil
}

// me only fetchs user's `id`
func me() (string, error) {
	b, err := reqest(http.MethodGet, provideURL("/me"), nil, func(s int) bool { return s == http.StatusOK })
	if err != nil {
		return "", err
	}

	data := &struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}{}
	if err = json.Unmarshal(b, data); err != nil {
		return "", err
	}
	return data.Data.ID, nil
}

type payload struct {
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Format       string   `json:"contentFormat"`
	Tags         []string `json:"tags,omitempty"`
	Status       string   `json:"publishStatus,omitempty"`
	License      string   `json:"license,omitempty"`
	CanonicalURL string   `json:"canonicalUrl,omitempty"`
}

func (p *payload) post() ([]byte, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	b, err = reqest(http.MethodPost, provideURL(fmt.Sprintf("/users/%s/posts", uid)), bytes.NewReader(b), func(s int) bool { return s == http.StatusCreated })
	if err != nil {
		return nil, err
	}
	return b, nil
}

func reqest(method, url string, payload io.Reader, ok func(status int) bool) ([]byte, error) {
	r, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	r.Header.Add("Content-Type", "application/json; charset=utf-8")
	r.Header.Add("Accept", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if !ok(resp.StatusCode) {
		return nil, fmt.Errorf("unwanted status code %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
