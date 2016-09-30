package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/longkai/xiaolongtongxue.com/helper"
)

var (
	fs     http.Handler
	v      string
	cdn    string
	origin string
)

// Init static file handler
func initFS(_cdn, _origin, _v string) {
	src, dest := `assets`, filepath.Join(env.Repo, `assets`)
	log.Printf("cpAssets(%q, %q)", src, dest)
	go cpAssets(src, dest)
	fs = http.FileServer(http.Dir(env.Repo))
	http.Handle(`/assets/`, fs)
	cdn, origin, v = _cdn, _origin, _v
	if cdn != "" {
		prefix := `/cdn/`
		log.Printf("http.StripPrefix(%q) for CDN %s", prefix, cdn)
		http.Handle(prefix, http.StripPrefix(prefix, fs))
	}
}

func serveFile(w http.ResponseWriter, r *http.Request) { fs.ServeHTTP(w, r) }

func cpAssets(src, dest string) {
	// ensure dir
	for _, e := range helper.Dirents(src) {
		_src, _dest := filepath.Join(src, e.Name()), filepath.Join(dest, e.Name())
		if e.IsDir() {
			go cpAssets(_src, _dest)
		} else {
			go func() {
				if err := helper.Cp(_src, _dest); err != nil {
					log.Print(err)
				}
			}()
		}
	}
}

// give assets a verion in its file name(e.g. for escape cdn cache)
func revAsset(name string) string {
	if cdn == "" || v == "" {
		return name
	}
	// TODO: should we plus image source?
	if strings.HasSuffix(name, ".js") || strings.HasSuffix(name, ".css") {
		i := strings.LastIndexByte(name, '.')
		return fmt.Sprintf("%s-%s%s", name[:i], v, name[i:])
	}
	return name
}

// TransformCDN to cdn href
func TransformCDN(href string) template.URL {
	if cdn == "" {
		return template.URL(href)
	}
	return template.URL(cdn + revAsset(href))
}
