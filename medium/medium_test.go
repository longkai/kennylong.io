package medium

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMe(t *testing.T) {
	var (
		status int
		body   = []byte(`
{
  "data": {
    "id": "5303d74c64f66366f00cb9b2a94f3251bf5",
    "username": "majelbstoat",
    "name": "Jamie Talbot",
    "url": "https://medium.com/@majelbstoat",
    "imageUrl": "https://images.medium.com/0*fkfQiTzT7TlUGGyI.png"
  }
}`)
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(body)
	}))
	defer ts.Close()

	saved := provideURL
	defer func() { provideURL = saved }()
	provideURL = func(p string) string { return ts.URL + p }
	m := &Medium{}
	t.Run("OK", func(t *testing.T) {
		status = http.StatusOK
		wantID := "5303d74c64f66366f00cb9b2a94f3251bf5"
		id, err := me(m.token)
		if err != nil {
			t.Errorf("me() fail: %v", err)
		}
		if id != wantID {
			t.Errorf("me() got id %q, want %q", id, wantID)
		}
	})

	t.Run("StatusFail", func(t *testing.T) {
		status = http.StatusBadRequest
		if _, err := me(m.token); err == nil {
			t.Errorf("me(_) status %d, want %d", http.StatusOK, http.StatusBadRequest)
		}
	})

	t.Run("MalformBody", func(t *testing.T) {
		status = http.StatusOK
		body = []byte(`balabala`)
		if _, err := me(m.token); err == nil {
			t.Errorf("me(_) should fail")
		}
	})
}

func TestPost(t *testing.T) {
	var (
		status int
		body   = []byte(`balabala`)
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(body)
	}))
	defer ts.Close()

	saved1 := provideURL
	defer func() { provideURL = saved1 }()
	provideURL = func(p string) string { return ts.URL + p }

	m := &Medium{}

	t.Run("OK", func(t *testing.T) {
		status = http.StatusCreated
		b, err := new(payload).post(m.uid, m.token)
		if err != nil {
			t.Errorf("post() fail: %v", err)
		}
		if !bytes.Contains(b, body) {
			t.Errorf("post() got %s, want contains %s", b, body)
		}
	})

	t.Run("Fail", func(t *testing.T) {
		status = http.StatusOK
		_, err := new(payload).post(m.uid, m.token)
		if err == nil {
			t.Errorf("post() should fail with status %d", status)
		}
		if !strings.Contains(err.Error(), fmt.Sprintf("%d", status)) {
			t.Errorf("post() fail with status: %s, want %d", err.Error(), status)
		}
	})
}
