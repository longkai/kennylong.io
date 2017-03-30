package repo

import (
	"log"
	"path/filepath"
	"strings"
	"sync"
)

// Processor process a path, which could be a file or dir.
type Processor interface {
	Process(path string) Docs
}

// DocsProcessor render mark up documents as articles.
type DocsProcessor struct {
	pt       PathTransformer
	scanner  Scanner
	parser   Parser
	callback func(docs Docs)
}

// Process the given for documents, either using callback
// or return value to receiving the results, or both.
func (p *DocsProcessor) Process(path string) Docs {
	// Ensure path absolute.
	if path = p.pt.Abs(path); path == "" {
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

func (p *DocsProcessor) process(path string, ch chan<- Doc) {
	url := p.pt.URLPath(path)
	if url == "" || url == "/" {
		log.Printf("skip process path: %q", path)
		return
	}

	doc, err := p.parser.Parse(path)
	if err != nil {
		log.Printf("process %q fail: %v", path, err)
		return
	}

	doc.Path = path
	doc.URL = url
	// URL is case insensitive.
	doc.Background = strings.ToLower(doc.Background)
	if bg := doc.Background; bg != "" && !strings.HasPrefix(bg, "/") &&
		!strings.HasPrefix(bg, "https://") && !strings.HasPrefix(bg, "http://") {
		// Normalize relative links.
		doc.Background = filepath.Join(doc.URL, doc.Background)
	}

	ch <- doc
}
