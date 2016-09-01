package github

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type statusHandler int

func (s *statusHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(int(*s))
}

func TestMarkdown(t *testing.T) {
	var status statusHandler
	ts := httptest.NewServer(&status)
	defer ts.Close()

	tests := []struct {
		status int
		err    error
	}{
		{http.StatusOK, nil},
		{415, errors.New("415")},
		{301, errors.New("301")},
		{500, errors.New("500")},
	}

	saved1, saved2 := provideURL, provideToken
	defer func() {
		provideURL = saved1
		provideToken = saved2
	}()
	// stub
	provideURL = func(path string) string {
		return ts.URL
	}
	provideToken = func() string {
		return ""
	}
	for _, test := range tests {
		status = statusHandler(test.status)
		if _, err := Markdown(strings.NewReader("...")); err != test.err {
			if !strings.Contains(err.Error(), test.err.Error()) {
				t.Errorf("Markdown(...) = (_, %q), want Contains %q", err.Error(), test.err.Error())
			}
		}
	}
}
