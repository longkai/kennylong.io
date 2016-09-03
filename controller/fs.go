package controller

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/longkai/xiaolongtongxue.com/helper"
)

var (
	fs http.Handler
)

// Init static file handler
func initFS(cdnPrefix string) {
	src, dest := `assets`, filepath.Join(repo, `assets`)
	log.Printf("cpAssets(%q, %q)", src, dest)
	go cpAssets(src, dest)
	fs = http.FileServer(http.Dir(repo))
	http.Handle(`/assets/`, fs)
	if cdnPrefix != "" {
		log.Printf("http.StripPrefix(%q) for CDN", cdnPrefix)
		http.Handle(cdnPrefix, http.StripPrefix(cdnPrefix, fs))
	}
}

func cdn(w http.ResponseWriter, r *http.Request) { fs.ServeHTTP(w, r) }

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
