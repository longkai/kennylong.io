package controller

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/longkai/xiaolongtongxue.com/context"
	"github.com/longkai/xiaolongtongxue.com/github"
	"github.com/longkai/xiaolongtongxue.com/medium"
	"github.com/longkai/xiaolongtongxue.com/repo"
)

const (
	pageSize = 7
)

var (
	conf       context.Conf
	repository repo.Repo
	staticFS   http.Handler
	templs     *template.Template
)

// Ctrl main controller.
func Ctrl(_conf context.Conf) {
	conf = _conf

	repository = repo.NewRepo(conf.RepoDir, conf.SkipDirs, conf.GlobDocs,
		conf.Github.User, conf.Github.Repo, medium.NewMedium(conf))

	github.Init("/api/github/hook", conf.RepoDir, conf.Github.HookSecret,
		conf.Github.AccessToken, func(a, m, d []string) {
			repository.Batch(a, m, d)
		})

	templs = template.Must(template.New("templ").ParseGlob("templ/*"))

	// Global handler.
	http.HandleFunc("/", handle)
}

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/": // Home page.
		home(w, r)
	case "/list": // Pagination.
		list(w, r)
	default: // Otherwise it should be an entry or its static resources request.
		entry(w, r)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	v := repository.List("/", pageSize)
	data := &struct {
		List repo.Docs
		Meta interface{}
	}{v, conf.Meta}

	if err := templs.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func list(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("since")
	// If the key is not found, return the latest ones.
	v := repository.List(key, pageSize)
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

// entry bundles entry and its static resources into a same handler,
// since we use the same directory layout, for doc and its static files(images).
// If entry is not found in repository, fallback to its static resources.
func entry(w http.ResponseWriter, r *http.Request) {
	// Compatible with old URL scheme, i.e., `/a/b/title/` to `/a/b/title`.
	p := r.URL.Path
	if strings.HasSuffix(p, "/") {
		http.Redirect(w, r, p[:len(p)-1], http.StatusMovedPermanently)
		return
	}
	// User-defined redirection mapping.
	if predir, ok := conf.Redirects[p]; ok {
		log.Printf("redirect %q to %q", p, predir)
		http.Redirect(w, r, predir, http.StatusMovedPermanently)
		return
	}
	// Try document first.
	doc, err := repository.Get(p)
	switch e := err.(type) {
	case nil:
		data := &struct {
			A    repo.Doc
			Meta interface{}
		}{doc, conf.Meta}
		// w.Header().Add("Cache-Control", "max-age=7200, public")
		if err = templs.ExecuteTemplate(w, "entry.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case repo.NotFoundError:
		// If no doc found, fallback to static files.
		staticFS.ServeHTTP(w, r)
	default: // general error
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}
