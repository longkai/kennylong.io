package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/longkai/xiaolongtongxue.com/config"
	"github.com/longkai/xiaolongtongxue.com/github"
	"github.com/longkai/xiaolongtongxue.com/medium"
	"github.com/longkai/xiaolongtongxue.com/render"
)

const (
	pageSize = 7
)

var (
	env    *config.Configuration
	sakura render.Engine
	templs *template.Template
)

// Ctrl main controller.
func Ctrl() {
	env = config.Env
	sakura = render.NewSakura(env.Meta.CDN)
	sakura.Post(env.Repo)
	installTempls()
	initFS(env.Meta.CDN, env.Meta.Origin, env.Meta.V)

	github.Init(`/api/github/hook`, env.Repo, env.HookSecret, env.AccessToken, revalidate)
	if env.MediumToken != "" {
		medium.Init(env.MediumToken, env.Meta.Origin)
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/ls/", ls)
	for _, v := range config.Roots() {
		installHanlder(v)
	}
}

var installTempls = func() {
	templs = template.Must(template.New(`sakura`).Funcs(template.FuncMap{
		`cdn`:     TransformCDN,
		`bgImg`:   render.BgImg,
		`tags`:    render.Tags,
		`format`:  render.Format,
		`daysAgo`: render.DaysAgo,
	}).ParseGlob(`templ/*`))
}

var installHanlder = func(p string) {
	p = fmt.Sprintf("/%s/", p)
	log.Printf("mapping url %s*", p)
	http.HandleFunc(p, entry)
}

var revalidate = func(a, m, d []string) {
	for i := range a {
		p := filepath.Join(env.Repo, a[i])
		// check if new router(i.e., URL starts with /balalaba/...)
		if v := config.Root(p); v != "" {
			installHanlder(v)
		}
		// TODO: better handling path travel stuffs...
		if env.MediumToken != "" && !config.Ignored(p) && strings.HasSuffix(p, ".md") {
			// meidum only allow posting new stuff, no other editing allow right now...
			go func() {
				if err := medium.Post(p); err != nil {
					log.Printf("medium.Post(%q) fail: %v", p, err)
				}
			}()
		}
	}
	if err := sakura.Revalidate(a, m, d); err != nil {
		log.Printf("revalidate fail: %v", err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		http.NotFound(w, r)
		return
	}

	v, err := sakura.Ls("", pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := &struct {
		List interface{}
		Meta interface{}
	}{v, config.Env.Meta}
	if err := templs.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func entry(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.RequestURI, "/") {
		serveFile(w, r)
		return
	}

	v, err := sakura.Get(r.RequestURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data := &struct {
		A    interface{}
		Meta interface{}
	}{v, config.Env.Meta}
	if err = templs.ExecuteTemplate(w, "entry.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ls(w http.ResponseWriter, r *http.Request) {
	key := r.RequestURI[len("/ls"):]
	if len(key) <= 1 { // `/` is not allowed
		http.Error(w, fmt.Sprintf("RequestURI %q, last segment not found", r.RequestURI), http.StatusBadRequest)
		return
	}
	v, err := sakura.Ls(key, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}
