package medium

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/longkai/xiaolongtongxue.com/context"
	"github.com/longkai/xiaolongtongxue.com/helper"
	"github.com/longkai/xiaolongtongxue.com/repo"
)

var provideURL = func(path string) string {
	return "https://api.medium.com/v1" + path
}

// Medium the medium writing platform.
type Medium struct {
	uid, token string
	origin     string
	renderer   repo.Renderer
	ready      chan struct{}
}

// NewMedium create a medium representation.
func NewMedium(conf context.Conf) *Medium {
	origin := conf.Meta.Origin
	m := &Medium{
		token:  conf.MediumToken,
		origin: origin,
		ready:  make(chan struct{}),
		renderer: &repo.GithubRenderer{
			User:       conf.Github.User,
			Repo:       conf.Github.Repo,
			Dir:        repo.Dir(conf.RepoDir),
			StripTitle: false,
			URLTransformer: func(str string) string {
				return fmt.Sprintf("%s/%s", origin, str)
			},
		},
	}
	go m.fetchUID()
	return m
}

// Visit repost the newly docs to the medium.
func (m *Medium) Visit(docs repo.Docs, cookie map[int]interface{}) {
	// Since medium only allows creating post, nothing we can do about editing.
	if v, ok := cookie[repo.Adds]; ok {
		if adds, ok := v.([]string); ok {
			sort.Strings(adds)
			for _, doc := range docs {
				i := sort.SearchStrings(adds, doc.Path)
				if i >= 0 && i < len(adds) {
					go m.Post(doc)
				}
			}
		}
	}
}

// It must be called before any othter functions in this package.
func (m *Medium) fetchUID() {
	val, err := helper.Try(3, func() (interface{}, error) { return me(m.token) })
	if err != nil {
		log.Printf("medium.me() fail: %v", err)
	} else {
		m.uid = val.(string)
		log.Printf("medium uid %s", m.uid)
	}
	close(m.ready)
}

// Post from the given path
func (m *Medium) Post(doc repo.Doc) error {
	<-m.ready
	if m.uid == "" {
		return fmt.Errorf("medium.Post(%q) without uid, abort posting", doc.Path)
	}

	html, err := m.renderer.Render(doc)
	if err != nil {
		return err
	}

	// Always `html` and `public` since you have publish to your site.
	p := &payload{Title: doc.Title, Tags: doc.Tags, License: doc.License, Format: "html", Status: "public"}
	p.CanonicalURL = m.origin + "/" + doc.URL
	p.Content = string(html)

	if _, err := p.post(m.uid, m.token); err != nil {
		return err
	}
	return nil
}

// me only fetchs user's `id`
func me(token string) (string, error) {
	b, err := reqest(http.MethodGet, provideURL("/me"), token, nil, func(s int) bool { return s == http.StatusOK })
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

func (p *payload) post(uid, token string) ([]byte, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	b, err = reqest(http.MethodPost, provideURL(fmt.Sprintf("/users/%s/posts", uid)), token, bytes.NewReader(b), func(s int) bool { return s == http.StatusCreated })
	if err != nil {
		return nil, err
	}
	return b, nil
}

func reqest(method, url, token string, payload io.Reader, ok func(status int) bool) ([]byte, error) {
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
