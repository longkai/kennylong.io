package github

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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

	saved := provideURL
	defer func() { provideURL = saved }()
	// stub
	provideURL = func(path string) string {
		return ts.URL
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

func TestHook(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(hook))
	defer ts.Close()

	setup := func() (*httptest.ResponseRecorder, *http.Request) {
		return httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, ts.URL, nil)
	}

	t.Run("BadSignature", func(t *testing.T) {
		w, r := setup()
		hook(w, r)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("hook(...) with Signature `blah`, want %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("GoodSignature", func(t *testing.T) {
		saved1, saved2 := handleHook, handleSecurity
		defer func() {
			handleHook = saved1
			handleSecurity = saved2
		}()
		// stub
		handleSecurity = func(reader io.Reader, signature string) error { return nil }

		// Note: concurrency here, so we need to wait for it to be finished
		var wg sync.WaitGroup

		called := false
		handleHook = func() {
			called = true
			wg.Done()
		}

		wg.Add(1)
		w, r := setup()
		r.Header.Add(`X-GitHub-Event`, `push`)
		hook(w, r)
		wg.Wait()

		if w.Code != http.StatusOK {
			t.Errorf("status %d, want %d", w.Code, http.StatusOK)
		}

		if !strings.Contains(w.Body.String(), "thx") {
			t.Errorf("strings.Contains(%q, 'thx') = false", w.Body.String())
		}

		if !called {
			t.Errorf("handleHook() not called")
		}
	})
}
