// traversal travel the markdown dirs and turn them into htmls
package render

import (
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/longkai/xiaolongtongxue.com/env"
	"github.com/longkai/xiaolongtongxue.com/github"
)

var (
	dirSema    = make(chan struct{}, 20) // max concurrent dir travel routine
	renderSema = make(chan struct{}, 5)  // max render routine
	entryTempl = template.Must(template.New("entry.html").Funcs(template.FuncMap{
		"format":   Format,
		"tags":     Tags,
		"hasColor": HasColor,
		"hasImage": HasImage,
	}).ParseFiles(env.Template+"/entry.html", env.Template+"/include.html"))
)

// perform a recursive directory walking, render all *.md to html(if not skipped or kept) to the same directory layout.
func Traversal(root string) Articles {
	metas := make(chan markdownMeta)
	var n sync.WaitGroup
	n.Add(1)
	doTraversal(root, &n, metas)
	go func() {
		n.Wait()
		close(metas)
	}()
	a := Articles{}
	for m := range metas {
		a = append(a, m)
	}
	sort.Sort(a)
	return a
}

func doTraversal(dir string, n *sync.WaitGroup, metas chan<- markdownMeta) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		fname := filepath.Join(dir, entry.Name())
		switch {
		case strings.HasPrefix(entry.Name(), "."):
			// dot file, ignore
		case entry.IsDir():
			if !env.Ignored(fname) {
				// dir, dive into
				n.Add(1)
				go doTraversal(fname, n, metas)
			}
		case strings.HasSuffix(entry.Name(), ".md"):
			if !env.Ignored(fname) {
				// a .md file, render it
				n.Add(1)
				go doRender(fname, n, metas)
			}
		default:
			// static file, copy it
			copyFile(fname, filepath.Join(env.Gen, fname[len(env.Config().ArticleRepo):]))
		}
	}
}

func doRender(fname string, n *sync.WaitGroup, metas chan<- markdownMeta) {
	defer n.Done()
	// acquire token
	renderSema <- struct{}{}
	defer func() {
		// release token
		<-renderSema
	}()
	m, err := newMarkdown(fname)
	if err != nil {
		// if anything fail, just skip it, same below
		log.Printf("render md %s fail, %v\n", fname, err)
		return
	}
	switch m.RenderOption {
	default:
		log.Printf("render option of %s is not valid.\n", fname)
		return
	case skip:
		log.Printf("%s is skipped.\n", fname)
		return
	case def:
	case keep:
		// see below
	}
	// save the file as index.html
	dest := filepath.Join(env.Gen, m.Id+"/index.html")
	err = ensureDir(dest)
	f, err := os.Create(dest)
	if err != nil {
		log.Printf("create dest file %s fail, %v\n", dest, err)
		return
	}

	defer func() {
		f.Close()
		m.Text, m.Html = "", "" // gc?
	}()

	// call github api
	b, err := github.Markdown(m.Text)
	if err != nil {
		log.Printf("render md %s fail, %v\n", fname, err)
		return
	}
	// write template to html
	m.Html = template.HTML(b)
	if err = entryTempl.Execute(f, m); err != nil {
		log.Printf("write %s entry template fail, %v\n", fname, err)
		return
	}

	// fmt.Printf("%s -> %s\n", fname, dest)
	m.Id = strings.Replace(m.Id, string(filepath.Separator), "/", -1) // adjust web path for windows...
	switch m.RenderOption {
	case def:
		metas <- m.markdownMeta
	case keep:
		log.Printf("%s is kept, will not appear in the article list.\n", fname)
	}
}

// dirents returns the entries of dir
// if nothing found, return nil
func dirents(dir string) []os.FileInfo {
	// acquire token
	dirSema <- struct{}{}
	defer func() {
		// release token
		<-dirSema
	}()

	f, err := os.Open(dir)
	if err != nil {
		log.Printf("open %s fail, %v\n", dir, err)
		return nil
	}

	defer f.Close()
	entries, err := f.Readdir(0) // 0 -> reads no limit
	if err != nil {
		log.Printf("Readdir fail, %v\n", err)
	}
	return entries
}

func copyFile(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	ensureDir(dest)
	f2, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f2.Close()
	defer f.Close()
	_, err = io.Copy(f2, f)
	return err
}

func ensureDir(dir string) error {
	dir = filepath.Dir(dir)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0744)
	}
	return err
}
