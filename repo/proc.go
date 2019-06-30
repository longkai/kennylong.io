package repo

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Processor process a path, which could be a file or dir.
type Processor interface {
	Process(paths ...string) Docs
}

// DocProcessor render mark up documents as articles.
type DocProcessor struct {
	dir      Dir
	scanner  Scanner
	parser   Parser
	callback func(docs Docs)
}

func rmExt(file string) string {
	f := filepath.Ext(file)
	index := strings.LastIndex(file, f)
	return file[:index]
}

// filter prefers md over org if a dir contains multiple valid docs.
func filter(files []string) []string {
	m := make(map[string]string)
	for _, file := range files {
		f := rmExt(file)
		if _, ok := m[f]; ok {
			if filepath.Ext(file) == ".md" {
				m[f] = file
			}
		} else {
			m[f] = file
		}
	}
	var res []string
	for _, v := range m {
		res = append(res, v)
	}
	return res
}

// Process the given for documents, either using callback
// or return value to receiving the results, or both.
func (p *DocProcessor) Process(paths ...string) Docs {
	var files []string
	for _, path := range paths {
		// Ensure absolute path.
		if path = p.dir.Abs(path); path == "" {
			log.Printf("skip outside path: %q", path)
			continue
		}
		files = append(files, filter(p.scanner.Scan(path))...)
	}

	var wg sync.WaitGroup
	ch := make(chan Doc)
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			p.process(file, ch)
		}(file)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var docs Docs
	for doc := range ch {
		docs = append(docs, doc)
	}

	// Send result back to receiver.
	if p.callback != nil {
		p.callback(docs)
	}
	return docs
}

func (p *DocProcessor) process(path string, ch chan<- Doc) {
	url := p.dir.URLPath(path)
	if url == "" {
		log.Printf("skip process path: %q", path)
		return
	}

	doc, err := p.parser.Parse(path)
	if err != nil {
		log.Printf("process %q fail: %v", path, err)
		return
	}

	doc.Path, doc.URL = path, url
	if bg := doc.Background; bg != "" && !absURLRegex.MatchString(bg) {
		// Complete relative links.
		doc.Background = filepath.Join(doc.URL, bg)
	}
	ch <- doc
}

// Anything that begins with a path relative links,
// e.g., no scheme(include relative scheme), no start with '/'.
var absURLRegex = regexp.MustCompile(`(^\w+:\/\/)|(^\/{1,2})`)
