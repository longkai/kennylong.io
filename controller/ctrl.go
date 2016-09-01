package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/longkai/xiaolongtongxue.com/config"
	"github.com/longkai/xiaolongtongxue.com/render"
)

const (
	pageSize = 7
)

var (
	templ    = `templ` // templates location
	homeTmpl = template.Must(template.New("index.html").Funcs(template.FuncMap{
		"daysAgo":  render.DaysAgo,
		"tags":     render.Tags,
		"hasColor": render.HasColor,
		"hasImage": render.HasImage,
		"relImage": render.IsRelImage,
	}).ParseFiles(templ+"/index.html", templ+"/include.html"))

	entryTempl = template.Must(template.New("entry.html").Funcs(template.FuncMap{
		"format":   render.Format,
		"tags":     render.Tags,
		"hasColor": render.HasColor,
		"hasImage": render.HasImage,
	}).ParseFiles(templ+"/entry.html", templ+"/include.html"))

	sakura   render.Engine
	staticFs http.Handler
)

// Ctrl main controller.
func Ctrl() {
	sakura = render.NewSakura()
	staticFs = http.FileServer(http.Dir(config.Env.ArticleRepo))
	sakura.Post(config.Env.ArticleRepo)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", home)
	http.HandleFunc("/ls/", ls)
	for _, dir := range config.Env.PublishDirs {
		http.Handle(fmt.Sprintf("/%s/", dir), http.HandlerFunc(entry))
	}
}

func home(w http.ResponseWriter, req *http.Request) {
	v, err := sakura.Ls("", pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := homeTmpl.Execute(w, v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func entry(w http.ResponseWriter, req *http.Request) {
	if !strings.HasSuffix(req.RequestURI, "/") {
		staticFs.ServeHTTP(w, req)
		return
	}

	v, err := sakura.Get(req.RequestURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = entryTempl.Execute(w, v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ls(w http.ResponseWriter, req *http.Request) {
	key := req.RequestURI[len("/ls"):]
	if len(key) <= 1 { // `/` is not allowed
		http.Error(w, fmt.Sprintf("RequestURI %q, last segment not found", req.RequestURI), http.StatusBadRequest)
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
