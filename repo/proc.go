package repo

import (
	"log"
	"path/filepath"
	"regexp"
	"sync"
)

// Processor process a path, which could be a file or dir.
type Processor interface {
	Process(path string) Docs
}

// DocProcessor render mark up documents as articles.
type DocProcessor struct {
	dir      Dir
	scanner  Scanner
	parser   Parser
	callback func(docs Docs)
}

// Process the given for documents, either using callback
// or return value to receiving the results, or both.
func (p *DocProcessor) Process(path string) Docs {
	// Ensure path absolute.
	if path = p.dir.Abs(path); path == "" {
		log.Printf("skip outside path: %q", path)
		return nil
	}

	files := p.scanner.Scan(path)

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
	if bg := doc.Background; bg != "" && !absURLRegex.MatchString(doc.Background) {
		// Complete relative links.
		doc.Background = filepath.Join(doc.URL, bg)
	}
	ch <- doc
}

// Anything that begins with a path relative links,
// e.g., no scheme(include relative scheme), no start with '/'.
var absURLRegex = regexp.MustCompile(`(^\w+:\/\/)|(^\/{1,2})`)
