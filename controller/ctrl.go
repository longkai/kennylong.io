package controller

import (
	"encoding/json"
	"html/template"
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
		conf.Github.AccessToken, func(a, m, d []string) { repository.Batch(a, m, d) })

	templs = template.Must(template.New("templ").ParseGlob("templ/*"))
	staticFS = http.FileServer(http.Dir(conf.RepoDir))

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/", handle)
	http.HandleFunc("/ls/", ls)
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		v := repository.List("", pageSize)
		data := &struct {
			List repo.Docs
			Meta interface{}
		}{v, conf.Meta}
		if err := templs.ExecuteTemplate(w, "index.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// Otherwise it should be a entry request.
	entry(w, r)
}

func entry(w http.ResponseWriter, r *http.Request) {
	// Compatible with old URL scheme, i.e., `/a/b/title/` to `/a/b/title`.
	if p := r.URL.Path; strings.HasSuffix(p, "/") {
		http.Redirect(w, r, p[:len(p)-1], http.StatusMovedPermanently)
		return
	}
	// Try article first.
	doc, err := repository.Get(r.URL.Path)
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
	case repo.NotFound:
		// If no doc found, fallback to static files.
		staticFS.ServeHTTP(w, r)
	default: // general error
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func ls(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/ls"):]
	v := repository.List(key, pageSize)
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}
