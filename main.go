package main

import (
	"encoding/json"
	"fmt"
	"github.com/longkai/xiaolongtongxue.com/env"
	"github.com/longkai/xiaolongtongxue.com/render"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
)

const (
	MAX_PAGE_SIZE = 20
)

type articles []render.MarkdownMeta

var (
	homeTmpl = template.Must(template.New("index.html").Funcs(template.FuncMap{
		"daysAgo":  render.DaysAgo,
		"tags":     render.Tags,
		"hasColor": render.HasColor,
		"hasImage": render.HasImage,
	}).ParseFiles(env.Template + "/index.html"))

	staticFs = http.FileServer(http.Dir(env.GEN))
)

func main() {
	env.InitEnv("testing_env.json")
	var list articles = render.Traversal(env.Config().ArticleRepo)
	sort.Sort(list)
	http.HandleFunc("/", list.home)
	http.HandleFunc("/articles/inc", list.inc)
	http.HandleFunc("/pagination", list.pagination)
	http.HandleFunc("/len", list.len)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	// article
	http.Handle("/articles/", http.HandlerFunc(article))
	fmt.Printf("\nHappy hackcing :)\n")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func article(resp http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.Path)
	staticFs.ServeHTTP(resp, req)
}

func (a articles) inc(resp http.ResponseWriter, req *http.Request) {
	if err := requiredMethod(http.MethodGet, req.Method); err != nil {
		http.Error(resp, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	id, inc := "", 0 // default to the zero index :)
	id = req.URL.Query().Get("me")
	inc, _ = strconv.Atoi(req.URL.Query().Get("inc"))

	// TODO: seq search, since it's not large
	i := -1
	for j, m := range a {
		if m.Id == id {
			i = j
			break
		}
	}

	if i != -1 {
		i := i + inc
		if 0 <= i && i < len(a) {
			// found it
			if b, err := json.Marshal(&a[i]); err != nil {
				log.Println(err)
				http.Error(resp, "", http.StatusInternalServerError)
			} else {
				resp.Header().Add("Content-Type", "application/json; charset=utf-8")
				fmt.Fprintf(resp, "%s", b)
			}
			return
		}
	}
	// otherwise, not found
	http.Error(resp, "", http.StatusNotFound)
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
	// lt max size
	if size > MAX_PAGE_SIZE {
		http.Error(resp, fmt.Sprintf("size required less than %d, got %d", MAX_PAGE_SIZE, size), http.StatusBadRequest)
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
	if req.URL.Path != "/" {
		http.Error(resp, "404 page not found", http.StatusNotFound)
		return
	}

	if err := requiredMethod(http.MethodGet, req.Method); err != nil {
		http.Error(resp, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	if err := homeTmpl.Execute(resp, a); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}
}

func (a articles) render(resp http.ResponseWriter) {
	b, err := json.Marshal(a)
	if err != nil {
		log.Println(err.Error()) // log locally
		http.Error(resp, "", http.StatusInternalServerError)
		return
	}

	resp.Header().Add("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(resp, "%s", b)
}

func requiredMethod(required, got string) error {
	if required != got {
		return fmt.Errorf("unsupport method, required %s, got %s\n", required, got)
	}
	return nil
}

func mustNature(n int, key string) error {
	if n < 0 {
		return fmt.Errorf("%q must be nature number, got %d.", key, n)
	}
	return nil
}

func (a articles) Len() int { return len(a) }

func (a articles) Less(i, j int) bool { return a[i].Date.After(a[j].Date) }

func (a articles) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
