package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	if req.key == "" {
		// from the max down to size
		it := s.list.SeekToLast()
		defer it.Close()
		l = append(l, it.Value())
		for i := 1; i < req.size && it.Previous(); i++ {
			l = append(l, it.Value())
		}
		go func() {
			req.resp <- l
		}()
		return
	}

	index, ok := s.index[req.key]
	if !ok {
		go func() {
			req.resp <- fmt.Errorf("key %q not found", req.key)
		}()
		return
	}
	it := s.list.Seek(index)
	defer it.Close()
	for i := 0; i < req.size && it.Previous(); i++ {
		l = append(l, it.Value())
	}
	// deliver
	go func() {
		req.resp <- l
	}()
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
	it := s.list.Seek(index)
	defer it.Close()
	// Note it's time asc order
	if it.Previous() {
		next = it.Value().(*Meta).ID
		it.Next() // go back
	}
	if it.Next() {
		prev = it.Value().(*Meta).ID
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

func (s *Sakura) travel(dir string) {
	dirSema <- struct{}{}
	defer func() { <-dirSema }()

	for _, e := range dirents(dir) {
		if s.Fun(dir, e) {
			go s.Meet(filepath.Join(dir, e.Name()))
		}
	}
}

// Fun is it?
func (s *Sakura) Fun(place string, sth os.FileInfo) bool {
	name := sth.Name()
	switch {
	default:
		// normal files, no actions done here
		return false
	case strings.HasPrefix(name, "."):
		// ignore hidden stuffs
		return false
	case sth.IsDir():
		// have no idea, so we need to look agian
		go s.travel(filepath.Join(place, name))
		return false
	case strings.HasSuffix(strings.ToLower(name), ".md"):
		// interesting :)
		return true
	}
}

// Meet sth. interesting.
func (s *Sakura) Meet(sth string) {
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
	s.requests.post <- m
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
	go s.travel(dir)
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
		Render: func(in io.Reader) (interface{}, error) {
			return github.Markdown(in)
		},
	}
	go s.loop()
	return s
}

// dirSema is a counting demaphore for limiting concurrency in dirents.
var dirSema = make(chan struct{}, 20)

// dirents lists the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Print(err)
		return nil
	}
	return entries
}
