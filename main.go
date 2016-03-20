package main

import (
	"encoding/json"
	"fmt"
	"github.com/longkai/xiaolongtongxue.com/env"
	"github.com/longkai/xiaolongtongxue.com/render"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type articles []render.MarkdownMeta

func main() {
	env.InitEnv("testing_env.json")
	var list articles = render.Traversal(env.Config().ArticleRepo)
	sort.Sort(list)
	http.HandleFunc("/", list.home)
	http.HandleFunc("/pagination", list.pagination)
	http.HandleFunc("/len", list.len)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

// this should be used rarely
func (a articles) len(resp http.ResponseWriter, req *http.Request) {
	if err := requiredMethod(http.MethodGet, req.Method); err != nil {
		http.Error(resp, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(resp, "%d", len(a))
}

func (a articles) pagination(resp http.ResponseWriter, req *http.Request) {
	if err := requiredMethod(http.MethodGet, req.Method); err != nil {
		http.Error(resp, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	p, size := 0, 7                               // default page and size per page
	p, _ = strconv.Atoi(req.URL.Query().Get("p")) // don' t care
	if str := req.URL.Query().Get("size"); str != "" {
		tmp, err := strconv.Atoi(str)
		if err == nil {
			size = tmp
		}
	}
	// ensure >= 0
	if err := mustNature(p, "p"); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	if err := mustNature(size, "size"); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	result := a.offset(p*size, size)
	result.render(resp)
}

func (a articles) offset(start, offset int) articles {
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

func (a articles) home(resp http.ResponseWriter, req *http.Request) {
	if req.URL.String() != "/" {
		http.Error(resp, "", http.StatusNotFound)
		return
	}

	if err := requiredMethod(http.MethodGet, req.Method); err != nil {
		http.Error(resp, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	a.render(resp)
}

func (a articles) render(resp http.ResponseWriter) {
	b, err := json.Marshal(a)
	if err != nil {
		http.Error(resp, "internal error", http.StatusInternalServerError)
		return
	}

	resp.Header().Add("Content-Type", "application/json;charset=utf-8")
	fmt.Fprintf(resp, "%s", b)
}

func requiredMethod(required, got string) error {
	if required != got {
		return fmt.Errorf("unsupport method, required %s, got %s\n", required, got)
	}
	return nil
}

func print(s []render.MarkdownMeta) {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
}

func mustNature(n int, key string) error {
	if n < 0 {
		return fmt.Errorf("%q must be nature number, got %d.", key, n)
	}
	return nil
}

func (a articles) Len() int { return len(a) }

func (a articles) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

func (a articles) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
