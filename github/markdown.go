package github

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/longkai/xiaolongtongxue.com/config"
)

const (
	endpoint = "https://api.github.com"
)

var provideURL = func(path string) string {
	return endpoint + path
}

var provideToken = func() string {
	return fmt.Sprintf("token %s", config.Env.AccessToken)
}

// Markdown makrdownify plain text to html with Github API.
func Markdown(in io.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, provideURL("/markdown/raw"), in)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")
	req.Header.Add("Authorization", provideToken())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Markdown API StatusCode %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
