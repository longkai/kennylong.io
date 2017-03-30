package repo

import (
	"html/template"
	"time"
)

// NotFound indicates the doc is not existed.
type NotFound string

func (s NotFound) Error() string { return "404 doc not found: " + string(s) }

// Doc representation of the document.
type Doc struct {
	Body       template.HTML `json:"-" yaml:"-"`
	Path       string        `json:"-" yaml:"-"`
	URL        string        `json:"url,omitempty" yaml:"-"`
	Title      string        `json:"title" yaml:",omitempty"`
	Tags       []string      `json:"tags"`
	Date       time.Time     `json:"date"`
	Weather    string        `json:"weather"`
	Summary    string        `json:"summary"`
	Location   string        `json:"location"`
	Background string        `json:"background"`
	License    string        `json:"license"`
	// hide from the list, but still can get with URL
	Hide  bool   `json:"hide"`
	Older string `json:"older,omitempty"`
	Newer string `json:"newer,omitempty"`
}

// Docs a list of articles. Maybe a BST is faster,
// an list is fast enough for such a tiny program, however.
type Docs []Doc

// Len _
func (d Docs) Len() int { return len(d) }

// Less order by time Dec.
func (d Docs) Less(i, j int) bool { return d[i].Date.After(d[j].Date) }

// Swap _
func (d Docs) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

func (d Docs) filterHidden(doc Doc) bool { return doc.Hide }

func (d Docs) travel(begin int, take int, inc bool, filter func(doc Doc) bool) Docs {
	var k = 0
	var res = Docs{}
	for i := begin; i < len(d) && i >= 0 && k < take; {
		if doc := d[i]; !filter(doc) {
			res = append(res, doc)
			k++
		}
		if inc {
			i++
		} else {
			i--
		}
	}
	return res
}

func (d Docs) matchFirst(match func(doc Doc) bool) int {
	for i, v := range d {
		if match(v) {
			return i
		}
	}
	return -1
}
