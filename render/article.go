package render

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Articles []MarkdownMeta

func (a Articles) Offset(start, offset int) Articles {
	if start > len(a) {
		// no more
		return a[len(a):]
	}
	if start+offset > len(a) {
		// not enough
		return a[start:]
	}
	return a[start : start+offset]
}

func (a Articles) Render(resp http.ResponseWriter) {
	b, err := json.Marshal(a)
	if err != nil {
		log.Println(err.Error()) // log locally
		http.Error(resp, "", http.StatusInternalServerError)
		return
	}

	resp.Header().Add("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(resp, "%s", b)
}

func (a Articles) Len() int { return len(a) }

func (a Articles) Less(i, j int) bool { return a[i].Date.After(a[j].Date) }

func (a Articles) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
