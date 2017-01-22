package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"path/filepath"

	"github.com/longkai/xiaolongtongxue.com/config"
	"github.com/longkai/xiaolongtongxue.com/github"
	"github.com/ryszard/goskiplist/skiplist"
)

type timestamp int64

// LessThan using timestamp from meta as key for skiplist sort.
func (t timestamp) LessThan(other skiplist.Ordered) bool { return t < timestamp(other.(timestamp)) }

type render func(id string, in io.Reader) (interface{}, error)

type entry struct {
	//val   interface{}
	err   error
	ready chan struct{}
}

func (e *entry) call(id string, in io.Reader, f render, callback func(b interface{})) {
	var (
		retry = 0
		val   interface{}
	)
	for retry < 3 { // max 3 retries
		val, e.err = f(id, in)
		if e.err == nil {
			// success
			callback(val)
			break
		}
		retry++
	}
	close(e.ready)
}

type request struct {
	key  string
	resp chan interface{}
}

type listrequest struct {
	key  string
	size int
	resp chan interface{}
}

type requests struct {
	get  chan request
	post chan interface{}
	put  chan request
	del  chan request
	list chan listrequest
}

// Sakura render engine.
type Sakura struct {
	render
	Traveller
	requests
	list  *skiplist.SkipList
	index map[string]timestamp
	cache map[string]*entry
}

func (s *Sakura) loop() {
	for {
		select {
		case req := <-s.requests.list:
			s.ls(req)
		case req := <-s.requests.get:
			s.get(req)
		case v := <-s.requests.post:
			s.post(v)
		case req := <-s.requests.del:
			s.del(req)
		}
	}
}

func (s *Sakura) del(req request) {
	// sync
	delete(s.cache, req.key)
	ts := s.index[req.key]
	delete(s.index, req.key)
	_, ok := s.list.Delete(ts)
	log.Printf("del %q: %t, total %d", req.key, ok, len(s.index))
	// deliver
	go func() { req.resp <- ok }()
}

func (s *Sakura) post(v interface{}) {
	// sync it
	m := v.(*Meta)
	if _, ok := s.index[m.ID]; ok {
		log.Printf("duplicated key %q\n", m.ID)
		return
	}
	ts := timestamp(m.Date.UnixNano())
	s.index[m.ID] = ts
	s.list.Set(ts, m)
	log.Printf("Find %q, total %d\n", m.ID, s.list.Len())
}

func (s *Sakura) ls(req listrequest) {
	l := []interface{}{}
	f := func(it skiplist.Iterator, size int) {
		for i := 0; i < size; {
			v := it.Value()
			// drop those who are hide
			if !v.(*Meta).Hide {
				l = append(l, it.Value())
				i++
			}
			if !it.Previous() {
				break
			}
		}
	}

	deliver := func(v interface{}) { req.resp <- v }

	if req.key == "" { // the index page request
		// from the max down to size
		it := s.list.SeekToLast()
		defer it.Close()
		f(it, req.size)
		go deliver(l)
		return
	}

	index, ok := s.index[req.key]
	if !ok {
		go deliver(fmt.Errorf("key %q not found", req.key))
		return
	}
	it := s.list.Seek(index)
	defer it.Close()
	// the key one has been delivered, drop it
	if it.Previous() {
		f(it, req.size)
	}
	// deliver
	go deliver(l)
}

func (s *Sakura) get(req request) {
	index, ok := s.index[req.key]
	if !ok {
		go func() {
			req.resp <- fmt.Errorf("key %q not found", req.key)
		}()
		return
	}
	v, _ := s.list.Get(index)
	var newer, older string
	// if it hide from the list, only show will directly http get access
	if !v.(*Meta).Hide {
		it := s.list.Seek(index)
		defer it.Close()
		var step int
		// Note it's time asc order
		for it.Previous() {
			step++
			v := it.Value()
			if !v.(*Meta).Hide {
				older = v.(*Meta).ID
				break
			}
		}
		// go back oringianl pos
		for step > 0 {
			it.Next()
			step--
		}
		for it.Next() {
			v := it.Value()
			if !v.(*Meta).Hide {
				newer = v.(*Meta).ID
				break
			}
		}
	}

	e := s.cache[req.key]
	m := v.(*Meta)
	if e == nil {
		// cache misses
		e = &entry{ready: make(chan struct{})}
		go e.call(m.ID, bytes.NewReader(m.Body.([]byte)), s.render, func(res interface{}) {
			m.Body = template.HTML(res.([]byte))
		})
		s.cache[req.key] = e
	}
	// deliver
	go func() {
		<-e.ready
		if e.err != nil {
			req.resp <- e.err
		} else {
			req.resp <- &Markdown{Meta: *m, Older: older, Newer: newer}
		}
	}()
}

// Ls markdown list.
func (s *Sakura) Ls(key string, size int) (interface{}, error) {
	resp := make(chan interface{})
	s.requests.list <- listrequest{key, size, resp}
	v := <-resp
	if err, ok := v.(error); ok {
		return nil, err
	}
	return v, nil
}

// Get markdown detail.
func (s *Sakura) Get(key string) (interface{}, error) {
	// log.Printf("Get(%q)", key)
	resp := make(chan interface{})
	s.requests.get <- request{key, resp}
	v := <-resp
	if err, ok := v.(error); ok {
		return nil, err
	}
	return v, nil
}

// Post markdowns for the given directory.
func (s *Sakura) Post(dir string) (interface{}, error) {
	log.Printf("Post(%q)", dir)
	go s.Traveller.Travel(dir)
	return nil, nil
}

// Put revalidate a markdown.
func (s *Sakura) Put(key string) (interface{}, error) {
	// if puts, del then post again since we don't know what has changes...
	log.Printf("Put(%q)", key)
	if _, err := s.Del(key); err != nil {
		log.Printf("del %q fail %v, exiting mod", key, err)
		return nil, err
	}
	// then post again
	s.Post(filepath.Join(config.Env.Repo, key))
	return nil, nil
}

// Del invalidate a markdown.
func (s *Sakura) Del(key string) (interface{}, error) {
	log.Printf("Del(%q)", key)
	resp := make(chan interface{})
	s.requests.del <- request{key, resp}
	v := <-resp
	if err, ok := v.(error); ok {
		return nil, err
	}
	return v, nil
}

// Revalidate the given entries.
func (s *Sakura) Revalidate(adds, mods, dels []string) error {
	log.Printf("Revalidate(%s, %s, %s)", adds, mods, dels)
	base := config.Env.Repo
	handle := func(a []string, f func(s string)) {
		for _, v := range a {
			sth := filepath.Join(base, v)
			// git only tracks files
			if s.Traveller.Into(sth) {
				f(sth)
			}
		}
	}

	handle(adds, func(sth string) { go s.Traveller.Meet(sth) })
	handle(mods, func(sth string) { go s.Put(parseID(sth)) })
	handle(dels, func(sth string) { go s.Del(parseID(sth)) })

	// always no error, since we can do nothing but logging
	return nil
}

// NewSakura sakura render engine with cdn support
func NewSakura(cdn string) Engine {
	s := &Sakura{
		cache: make(map[string]*entry),
		index: make(map[string]timestamp),
		list:  skiplist.New(),
		requests: requests{
			list: make(chan listrequest),
			del:  make(chan request),
			get:  make(chan request),
			post: make(chan interface{}),
			put:  make(chan request),
		},
		render: func(id string, in io.Reader) (interface{}, error) {
			b, err := github.Markdown(in)
			if err == nil && cdn != "" {
				return linkifyHTML(bytes.NewReader(b), []byte(cdn+id))
			}
			return b, err
		},
	}
	s.Traveller = &Hiker{func(v interface{}) { s.requests.post <- v }}
	go s.loop()
	return s
}
