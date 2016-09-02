package render

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/longkai/xiaolongtongxue.com/config"
	"github.com/longkai/xiaolongtongxue.com/helper"
)

// Traveller a traveller travels somewhere to meet sth. interesting.
type Traveller interface {
	// Travel a place to find someting interesting.
	Travel(place string)
	// Into if and only if sth is interesting.
	Into(sth string) bool
	// Meet meet with sth if it's really interesting.
	Meet(sth string)
}

// Hiker is a serious traveller.
type Hiker struct {
	callback func(interface{})
}

// Travel a place to find someting interesting.
func (h *Hiker) Travel(place string) {
	for _, e := range helper.Dirents(place) {
		sth := filepath.Join(place, e.Name())
		switch {
		case config.Ignored(sth):
		case e.IsDir():
			go h.Travel(sth)
		case h.Into(sth):
			go h.Meet(sth)
		}
	}
}

// Into if and only if sth is interesting.
func (h *Hiker) Into(sth string) bool {
	// TODO: support more ext?
	return !config.Ignored(sth) && strings.HasSuffix(sth, `.md`)
}

// Meet with sth if it's really interesting the traveller will call you.
func (h *Hiker) Meet(sth string) {
	f, err := os.Open(sth)
	if err != nil {
		log.Printf("open %q fail: %v", sth, err)
		return
	}
	defer f.Close()

	m, err := parseMd(f)
	if err != nil {
		log.Printf("parse %q fail: %v", sth, err)
		return
	}
	m.ID = parseID(sth)

	h.callback(m) // let sb. know it's funny
}
