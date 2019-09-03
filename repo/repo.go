package repo

import (
	"html/template"
	"log"
	"sort"

	"github.com/longkai/xiaolongtongxue.com/helper"
)

// Repo the documents repository.
type Repo interface {
	List(since string, size int) Docs
	Get(path string) (Doc, error)
	Del(path string)
	Put(path string)
	Post(path string)
	Batch(adds, mods, dels []string)
}

type entry struct {
	ready chan struct{}
	val   template.HTML
	err   error
}

func (e *entry) call(doc Doc, r Renderer) {
	val, err := helper.Try(3, func() (interface{}, error) {
		return r.Render(doc)
	})
	e.val, e.err = val.(template.HTML), err
	close(e.ready)
}

type listReq struct {
	path string
	size int
	resp chan Docs
}

type batchReq struct {
	adds, mods, dels []string
}

type getReq struct {
	path string
	resp chan getResp
}

type getResp struct {
	doc Doc
	err error
}

type reqs struct {
	get   chan getReq
	post  chan Docs
	list  chan listReq
	batch chan batchReq
}

// DocRepo documents repository implements.
type DocRepo struct {
	reqs

	dir       Dir
	renderer  Renderer
	processor Processor
	visitors  []Visitor

	docs  Docs              // More read, less write.
	index map[string]int    // Fast lookup.
	cache map[string]*entry // Rendering cache.
}

func (r *DocRepo) loop() {
	for {
		select {
		case req := <-r.reqs.get:
			r.get(req)
		case doc := <-r.reqs.post:
			r.post(doc)
		case req := <-r.reqs.list:
			r.list(req)
		case req := <-r.reqs.batch:
			r.batch(req)
		}
	}
}

func (r *DocRepo) batch(req batchReq) {
	rm := func(paths []string) {
		for _, path := range paths {
			if _, ok := r.index[path]; ok {
				delete(r.index, path)
				delete(r.cache, path)
			}
		}
	}
	req.adds = filter(req.adds)
	req.mods = filter(req.mods)
	// Re-process adds and mods.
	if plen := len(req.adds) + len(req.mods); plen > 0 {
		combine := make([]string, plen)
		copy(combine, req.adds)
		copy(combine[len(req.adds):], req.mods)
		go func() {
			docs := r.processor.Process(combine...)
			// Post process.
			for _, v := range r.visitors {
				v.Visit(docs, map[int]interface{}{
					Adds: req.adds,
					Mods: req.mods,
				})
			}
		}()
	}
	// Deletions.
	rm(req.dels)
	// Modifications.
	rm(req.mods)
	// Rearrangement, strip dels and mods, the order still remains.
	idx := 0
	tmp := make(Docs, 0, len(r.docs))
	for _, doc := range r.docs {
		// Pick those who are not deleted or modified.
		if _, ok := r.index[doc.URL]; ok {
			tmp = append(tmp, doc)
			r.index[doc.URL] = idx
			idx++
		}
	}
	r.docs = tmp
	// Hence, the time complexity is O(n).
}

func (r *DocRepo) get(req getReq) {
	i, ok := r.index[req.path]
	if !ok {
		go func() { req.resp <- getResp{Doc{}, NotFoundError(req.path)} }()
		return
	}

	doc := r.docs[i]

	// A hidden doc has no newer/older navigation.
	if !doc.Hide {
		if docs := r.docs.travel(i+1, 1, true, r.docs.filterHidden); docs.Len() > 0 {
			doc.Older = docs[0].URL
		}
		if docs := r.docs.travel(i-1, 1, false, r.docs.filterHidden); docs.Len() > 0 {
			doc.Newer = docs[0].URL
		}
	}

	e := r.cache[req.path]
	if e == nil {
		// Cache misses.
		e = &entry{ready: make(chan struct{})}
		go e.call(doc, r.renderer)
		r.cache[req.path] = e
	}

	go func() {
		<-e.ready
		doc.Body = e.val
		req.resp <- getResp{doc, e.err}
	}()
}

func (r *DocRepo) post(docs Docs) {
	oldSize := r.docs.Len()
	// If a some of a newly doc has been existed already, replace them.
	// i.e., avoid duplication.
	for _, d := range docs {
		if i, ok := r.index[d.URL]; ok {
			log.Printf("replace %q with %q", r.docs[i].URL, d.URL)
			r.docs[i] = d          // Replace the old one.
			delete(r.cache, d.URL) // Clear its rendering cache, if any.
		} else {
			r.docs = append(r.docs, d) // Append the new one.
			// If there are multiple docs have a same URL,
			// the last one will be kept.
			r.index[d.URL] = r.docs.Len() - 1
		}
	}
	sort.Sort(r.docs)

	// Rebuild index.
	index := r.index
	for i, d := range r.docs {
		index[d.URL] = i
	}

	log.Printf("receive %d docs, len %d to %d", docs.Len(), oldSize, r.docs.Len())
}

func (r *DocRepo) list(req listReq) {
	i, ok := r.index[req.path]
	if !ok {
		i = 0 // If not found any match, start from 0.
	} else {
		i++ // Skip the current one.
	}

	res := r.docs.travel(i, req.size, true, r.docs.filterHidden)

	go func() { req.resp <- res }()
}

// List articles since a specific path, excluded.
func (r *DocRepo) List(since string, size int) Docs {
	resp := make(chan Docs)
	r.reqs.list <- listReq{since, size, resp}
	return <-resp
}

// Get a document for the path.
func (r *DocRepo) Get(path string) (Doc, error) {
	// Read only index, fast indexing without channel synchronization.
	// It's safe since lookup success, the go-routine will lookup again.
	// Hence, when lookup fail, maybe it's just removed or never exists.
	if _, ok := r.index[path]; !ok {
		return Doc{}, NotFoundError(path)
	}
	resp := make(chan getResp)
	r.reqs.get <- getReq{path, resp}
	v := <-resp
	return v.doc, v.err
}

// Del a document for the path.
func (r *DocRepo) Del(path string) {
	r.Batch(nil, nil, []string{path})
}

// Put revalidate a document.
func (r *DocRepo) Put(path string) {
	r.Batch(nil, []string{path}, nil)
}

// Post publish the path for documents.
// This method should be called when you start the application.
// Since it won't call the visitors let them do post process.
func (r *DocRepo) Post(path string) {
	log.Printf("post %s", path)
	r.processor.Process(path)
}

// Batch additions, modifications and deletions into a single request.
func (r *DocRepo) Batch(adds, mods, dels []string) {
	// Git only tracks files, hence, all of the slice are files path.

	// `git mv a b`: a deletion plus a addition, a and b is different.
	// `git rm`: deletion only.
	// `git add`: addition or modification.

	// Hence, the strategy is:
	// 1. deletion: just delete it from slice and rendering cache.
	// 2. modification: delete first then adding.
	// 3. addition: adding it.
	// 4. renaming: a delete and a addition.

	// What about a user modifies follows a renaming? A deletion and addition.

	// The key point: `adds`, `mods` and `dels` slice are distinct.
	// Therefore, order doesn't matter.

	log.Printf("Batch(%v, %v, %v)", filter(adds), filter(mods), filter(dels))

	r.reqs.batch <- batchReq{adds, mods, dels}
}

// NewRepo create a new article repository.
func NewRepo(repoDir string, skipDirs, globDocs []string,
	user, repo string, vistors ...Visitor) Repo {
	dir := Dir(repoDir)

	p := &DocProcessor{dir: dir}
	p.scanner = &DocScanner{
		dir:      dir,
		skipDirs: skipDirs,
		globDocs: globDocs,
	}
	p.parser = &DocParser{}

	r := new(DocRepo)
	r.dir = dir
	r.cache = make(map[string]*entry)
	r.index = make(map[string]int)
	r.processor = p
	r.visitors = vistors

	r.reqs = reqs{
		list:  make(chan listReq),
		get:   make(chan getReq),
		post:  make(chan Docs),
		batch: make(chan batchReq),
	}

	r.renderer = &newRender{
		wrap: NewRenderer(user, repo, dir),
	}

	// Receive result asynchronously.
	p.callback = func(docs Docs) { r.reqs.post <- docs }

	go r.loop()
	go r.Post(repoDir)
	return r
}
