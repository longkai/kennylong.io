package repo

import (
	"github.com/longkai/xiaolongtongxue.com/helper"
	"html/template"
	"log"
	"sort"
)

// Repo the articles repository.
type Repo interface {
	List(since string, size int) []Doc
	Get(path string) (Doc, error)
	Del(path string)
	Put(path string)
	Post(path string)
	Batch(adds, mods, dels []string)
}

// Visitor post render operation.
type Visitor interface {
	Visit(docs Docs)
}

type entry struct {
	ready chan struct{}
	val   template.HTML
	err   error
}

func (e *entry) call(path string, r Renderer) {
	val, err := helper.Retry(3, func() (interface{}, error) { return r.Render(path) })
	e.val, e.err = val.(template.HTML), err
	close(e.ready)
}

type listReq struct {
	key  string
	size int
	resp chan Docs
}

type getReq struct {
	path string
	resp chan getResp
}

type getResp struct {
	doc Doc
	err error
}

type requests struct {
	get  chan getReq
	post chan Docs
	put  chan string
	del  chan string
	list chan listReq
}

// Repository articles repository implements.
type Repository struct {
	requests
	visitors  []Visitor
	pt        PathTransformer
	renderer  Renderer
	processor Processor

	// Rendering cache.
	cache map[string]*entry
	docs  Docs
}

func (r *Repository) loop() {
	for {
		select {
		case req := <-r.requests.list:
			r.list(req)
		case req := <-r.requests.get:
			r.get(req)
		case path := <-r.requests.del:
			r.del(path)
		case path := <-r.requests.put:
			r.put(path)
		case article := <-r.requests.post:
			r.post(article)
		}
	}
}

func (r *Repository) list(req listReq) {
	var i = 0

	if req.key != "" {
		i = r.docs.matchFirst(func(doc Doc) bool {
			return doc.URL == req.key
		}) + 1 // Skip the current one.
		// If not found any match, start from 0.
	}

	res := r.docs.travel(i, req.size, true, r.docs.filterHidden)

	go func() { req.resp <- res }()
}

func (r *Repository) get(req getReq) {
	i := r.docs.matchFirst(func(doc Doc) bool {
		return doc.URL == req.path
	})
	if i < 0 {
		go func() { req.resp <- getResp{Doc{}, NotFound(req.path)} }()
		return
	}

	doc := r.docs[i]

	// A hidden doc has no newer/older navigation.
	if !doc.Hide {
		if docs := r.docs.travel(i+1, 1, true, r.docs.filterHidden); len(docs) > 0 {
			doc.Older = docs[0].URL
		}
		if docs := r.docs.travel(i-1, 1, false, r.docs.filterHidden); len(docs) > 0 {
			doc.Newer = docs[0].URL
		}
	}

	e := r.cache[req.path]
	if e == nil {
		// Cache misses.
		e = &entry{ready: make(chan struct{})}
		go e.call(doc.Path, r.renderer)
		r.cache[req.path] = e
	}

	go func() {
		<-e.ready
		doc.Body = e.val
		req.resp <- getResp{doc, e.err}
	}()
}

func (r *Repository) del(path string) {
	i := r.docs.matchFirst(func(doc Doc) bool {
		return doc.URL == path
	})
	if i > -1 {
		copy(r.docs[i:], r.docs[i+1:])
		r.docs = r.docs[:len(r.docs)-1]
		delete(r.cache, path)
	}
}

func (r *Repository) put(path string) {
	// Delete then repost.
	r.del(path)
	go r.processor.Process(path)
}

func (r *Repository) post(docs Docs) {
	r.docs = append(r.docs, docs...)
	sort.Sort(r.docs)
	log.Printf("receive %d new articles, total %d", len(docs), len(r.docs))
}

// List articles since a specific path, excluded.
func (r *Repository) List(since string, size int) []Doc {
	resp := make(chan Docs)
	r.requests.list <- listReq{since, size, resp}
	return <-resp
}

// Get a document for the path.
func (r *Repository) Get(path string) (Doc, error) {
	resp := make(chan getResp)
	r.requests.get <- getReq{path, resp}
	v := <-resp
	return v.doc, v.err
}

// Del a document for the path.
func (r *Repository) Del(path string) {
	log.Printf("delete: %s", path)
	r.requests.del <- path
}

// Put revalidate a document.
func (r *Repository) Put(path string) {
	log.Printf("revalidate: %s", path)
	r.requests.put <- path
}

// Post process the path for documents.
func (r *Repository) Post(path string) {
	log.Printf("post: %s", path)
	go r.processor.Process(path)
}

// Batch batch reqeust.
func (r *Repository) Batch(adds, mods, dels []string) {
	log.Printf("Batch(%v, %v, %v)", adds, mods, dels)
	// Git only tracks files, hence, all of the slice are files path.
	handle := func(a []string, f func(s string)) {
		for _, v := range a {
			f(v)
		}
	}

	handle(adds, func(f string) {
		go func() {
			docs := r.processor.Process(f)
			for _, v := range r.visitors {
				v.Visit(docs)
			}
		}()
	})
	handle(mods, func(f string) { go r.Put(r.pt.URLPath(f)) })
	handle(dels, func(f string) { go r.Del(r.pt.URLPath(f)) })
}

// NewRepo create a new article repository.
func NewRepo(repoDir string, skipDirs, globDocs []string, user, repo string, vistors ...Visitor) Repo {
	pt := NewPathTransformer(repoDir)

	p := &DocsProcessor{pt: pt}
	p.scanner = &DocScanner{
		skipDirs: skipDirs,
		globDocs: globDocs,
		pt:       pt,
	}
	p.parser = &DocParser{}

	r := new(Repository)
	r.pt = pt
	r.cache = make(map[string]*entry)
	r.processor = p
	r.visitors = vistors

	r.requests = requests{
		list: make(chan listReq),
		del:  make(chan string),
		put:  make(chan string),
		get:  make(chan getReq),
		post: make(chan Docs),
	}

	r.renderer = NewRenderer(user, repo, pt)

	p.callback = func(docs Docs) { r.requests.post <- docs }

	go r.loop()
	r.Post(repoDir)
	return r
}