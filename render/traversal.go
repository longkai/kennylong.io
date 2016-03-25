// traversal travel the markdown dirs and turn them into htmls
package render

import (
	"fmt"
	"github.com/longkai/xiaolongtongxue.com/env"
	"github.com/longkai/xiaolongtongxue.com/github"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

var (
	dirSema    = make(chan struct{}, 20) // max concurrent dir travel routine
	renderSema = make(chan struct{}, 5)  // max render routine
	entryTempl = template.Must(template.New("entry.html").Funcs(template.FuncMap{
		"daysAgo":  DaysAgo,
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
			copyFile(fname, filepath.Join(env.GEN, fname[len(env.Config().ArticleRepo):]))
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
		// if render fail, just skip it, same below
		fmt.Fprintf(os.Stderr, "render md %s fail, %v\n", fname, err)
		return
	}
	if m.RenderOption == skip {
		fmt.Printf("%s is skipped.\n", fname)
		return
	}
	b, err := github.Markdown(m.Text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render md %s fail, %v\n", fname, err)
		return
	}
	// save the file as index.html
	dest := filepath.Join(env.GEN, m.Id+"/index.html")
	err = ensureDir(dest)
	// fmt.Printf("%s -> %s\n", fname, dest)
	if f, err := os.Create(dest); err != nil {
		fmt.Fprintf(os.Stderr, "create dest file fail, %v\n", err)
	} else {
		// write the template
		m.Html = template.HTML(b)
		if err = entryTempl.Execute(f, m); err != nil {
			fmt.Fprintf(os.Stderr, "write entry template fail, %v\n", err)
		} else {
			switch m.RenderOption {
			case def:
				metas <- m.markdownMeta
			case keep:
				fmt.Printf("%s is kept, will not appear in the article list.", fname)
			}
		}
		m.Text, m.Html = "", "" // gc?
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
		fmt.Fprintf(os.Stderr, "open %s fail, %v\n", dir, err)
		return nil
	}

	defer f.Close()
	entries, err := f.Readdir(0) // 0 -> reads no limit
	if err != nil {
		fmt.Fprintf(os.Stdout, "Readdir fail, %v\n", err)
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
