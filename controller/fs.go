package controller

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/longkai/xiaolongtongxue.com/helper"
)

var (
	fs     http.Handler
	cdn    string
	domain string
)

// Init static file handler
func initFS(_cdn, _domain string) {
	src, dest := `assets`, filepath.Join(repo, `assets`)
	log.Printf("cpAssets(%q, %q)", src, dest)
	go cpAssets(src, dest)
	fs = http.FileServer(http.Dir(repo))
	http.Handle(`/assets/`, fs)
	cdn, domain = _cdn, _domain
	if cdn != "" {
		if domain == "" {
			log.Fatalf("CDN %q is enbaled, domain is empty", cdn)
		}
		prefix := `/cdn/`
		log.Printf("http.StripPrefix(%q) for CDN %s", prefix, cdn)
		http.Handle(prefix, http.StripPrefix(prefix, fs))
	}
}

func serveFile(w http.ResponseWriter, r *http.Request) { fs.ServeHTTP(w, r) }

// EscapeCDN if CDN is used, linkify those non-static(avoid CDN) links.
func EscapeCDN(url string) string {
	if cdn != "" {
		return "//" + domain + "/" + strings.TrimLeft(url, "/") // use the same protocol
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
