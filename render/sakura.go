package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"

	"github.com/longkai/xiaolongtongxue.com/github"
	"github.com/ryszard/goskiplist/skiplist"
)

// Render _
type Render func(in io.Reader) (interface{}, error)

type timestamp int64

// LessThan using timestamp from meta as key for skiplist sort.
func (t timestamp) LessThan(other skiplist.Ordered) bool {
	return t < timestamp(other.(timestamp))
}

type entry struct {
	val   interface{}
	err   error
	ready chan struct{}
}

func (e *entry) call(in io.Reader, f Render) {
	e.val, e.err = f(in)
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
	Render
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
		}
	}
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
		for i := 0; i < size && it.Previous(); {
			v := it.Value()
			// drop those who is hide
			if !v.(*Meta).Hide {
				l = append(l, it.Value())
				i++
			}
		}
	}

	deliver := func(v interface{}) { req.resp <- v }

	if req.key == "" {
		// from the max down to size
		it := s.list.SeekToLast()
		defer it.Close()
		l = append(l, it.Value())
		f(it, req.size-1)
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
	f(it, req.size)
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
	var prev, next string
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
				next = it.Value().(*Meta).ID
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
				prev = v.(*Meta).ID
			}
		}
	}

	e := s.cache[req.key]
	if e == nil {
		// cache misses
		e = &entry{ready: make(chan struct{})}
		go e.call(bytes.NewReader(v.(*Meta).body), s.Render)
		v.(*Meta).body = nil // clear unwanted data
		s.cache[req.key] = e
	}
	// deliver
	go func() {
		<-e.ready
		if e.err != nil {
			req.resp <- e.err
		} else {
			req.resp <- &Markdown{Meta: *v.(*Meta), Next: next, Prev: prev, Body: template.HTML(e.val.([]byte))} // TODO: cache the HTML better?
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
	go s.Traveller.Travel(dir)
	return nil, nil
}

// Put revalidate a markdown.
func (s *Sakura) Put(key string) (interface{}, error) {
	return nil, nil
}

// Del invalidate a markdown.
func (s Sakura) Del(key string) (interface{}, error) {
	return nil, nil
}

// NewSakura _
func NewSakura() Engine {
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
		Render: func(in io.Reader) (interface{}, error) { return github.Markdown(in) },
	}
	s.Traveller = &Hiker{func(v interface{}) { s.requests.post <- v }}
	go s.loop()
	return s
}
