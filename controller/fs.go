package controller

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/longkai/xiaolongtongxue.com/helper"
)

var (
	fs     http.Handler
	cdn    string
	origin string
)

// Init static file handler
func initFS(_cdn, _origin, v string) {
	src, dest := `assets`, filepath.Join(env.Repo, `assets`)
	log.Printf("cpAssets(%q, %q)", src, dest)
	go cpAssets(src, dest)
	fs = http.FileServer(http.Dir(env.Repo))
	http.Handle(`/assets/`, fs)
	cdn, origin = _cdn, _origin
	if cdn != "" {
		if origin == "" {
			log.Fatalf("CDN %q is enbaled, origin is empty", cdn)
		}
		if v != "" {
			v += "/" // for pretty URL
		}
		prefix := fmt.Sprintf("/cdn/%s", v) // plus a version code avoiding cdn hard cache...
		log.Printf("http.StripPrefix(%q) for CDN %s", prefix, cdn)
		http.Handle(prefix, http.StripPrefix(prefix, fs))
	}
}

func serveFile(w http.ResponseWriter, r *http.Request) { fs.ServeHTTP(w, r) }

// EscapeCDN if CDN is used, linkify those non-static(avoid CDN) links.
func EscapeCDN(url string) string {
	if cdn != "" {
		return origin + "/" + strings.TrimLeft(url, "/")
	}
	return url
}

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
