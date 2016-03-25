package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/longkai/xiaolongtongxue.com/env"
	"github.com/longkai/xiaolongtongxue.com/github"
	"github.com/longkai/xiaolongtongxue.com/render"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const (
	MAX_PAGE_SIZE     = 20
	DEFAULT_PAGE_SIZE = 7
)

var (
	homeTmpl = template.Must(template.New("index.html").Funcs(template.FuncMap{
		"daysAgo":  render.DaysAgo,
		"tags":     render.Tags,
		"hasColor": render.HasColor,
		"hasImage": render.HasImage,
		"relImage": render.IsRelImage,
	}).ParseFiles(env.Template+"/index.html", env.Template+"/include.html"))

	staticFs = http.FileServer(http.Dir(env.GEN))

	requests       = make(chan struct{})        // clients want our data
	responses      = make(chan render.Articles) // response to client with the data
	invalidate     = make(chan struct{})        // like android, calling this perform a redraw(re-render all articles)
	postInvalidate = make(chan render.Articles) // send the newly rendered articles back to looper
)

// looper confine all data for safety concurrency
// NOTE: we just simply return the data we hold, it's okay since it just a blog =.=, no write to the data itsefl now :)
func looper() {
	a := render.Traversal(env.Config().ArticleRepo)
	fmt.Printf("\nTotal article: %d, Happy hackcing :)\n", len(a))
	for {
		select {
		case <-invalidate:
			go doInvalidate()
		case <-requests:
			responses <- a
		case newly := <-postInvalidate:
			a = newly
		}
	}
}

func doInvalidate() {
	articles := render.Traversal(env.Config().ArticleRepo)
	log.Printf("hook -> pull -> reload %d articles :) \n", len(articles))
	postInvalidate <- articles
}

func main() {
	port := flag.Int("port", 1217, "http port number")
	conf := flag.String("conf", "testing_env.json", "config file path")
	flag.Parse()
	env.InitEnv(*conf)

	go looper()

	http.HandleFunc("/", home)
	http.HandleFunc("/articles/inc", inc)
	http.HandleFunc("/pagination", pagination)
	http.HandleFunc("/count", count)
	// api
	http.Handle("/api/", http.StripPrefix("/api", http.HandlerFunc(api)))
	// frontend static files
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	// gen articles
	for _, dir := range env.Config().PublishDirs {
		http.Handle(fmt.Sprintf("/%s/", dir), http.HandlerFunc(article))
	}
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func api(resp http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/github/hook": // github webhook
		github.Hook(resp, req, invalidate)
	}
}

func article(resp http.ResponseWriter, req *http.Request) {
	staticFs.ServeHTTP(resp, req)
}

func inc(resp http.ResponseWriter, req *http.Request) {
	if err := requiredMethod(http.MethodGet, req.Method); err != nil {
		http.Error(resp, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	id, inc := "", 0 // default to the zero index :)
	id = req.URL.Query().Get("me")
	inc, _ = strconv.Atoi(req.URL.Query().Get("inc"))

	// TODO: seq search, since it's not large
	requests <- struct{}{}
	a := <-responses
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
func count(resp http.ResponseWriter, req *http.Request) {
	if err := requiredMethod(http.MethodGet, req.Method); err != nil {
		http.Error(resp, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	requests <- struct{}{}
	a := <-responses
	fmt.Fprintf(resp, "%d", len(a))
}

func pagination(resp http.ResponseWriter, req *http.Request) {
	if err := requiredMethod(http.MethodGet, req.Method); err != nil {
		http.Error(resp, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	p, size := 0, DEFAULT_PAGE_SIZE               // default page and size per page
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

	requests <- struct{}{}
	a := <-responses
	result := a.Offset(p*size, size)
	result.Render(resp)
}

func home(resp http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.Error(resp, "404 page not found", http.StatusNotFound)
		return
	}

	if err := requiredMethod(http.MethodGet, req.Method); err != nil {
		http.Error(resp, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	requests <- struct{}{}
	a := <-responses
	var out render.Articles
	if len(a) < DEFAULT_PAGE_SIZE {
		out = a
	} else {
		out = a[:DEFAULT_PAGE_SIZE]
	}
	if err := homeTmpl.Execute(resp, out); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}
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
