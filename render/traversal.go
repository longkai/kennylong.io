// traversal travel the markdown dirs and turn them into htmls
package render

import (
	"fmt"
	"github.com/longkai/xiaolongtongxue.com/env"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	dirSema    = make(chan struct{}, 20) // max concurrent dir travel routine
	renderSema = make(chan struct{}, 5)  // max render routine
)

func Traversal(root string) []MarkdownMeta {
	metas := make(chan MarkdownMeta)
	var n sync.WaitGroup
	n.Add(1)
	doTraversal(root, &n, metas)
	go func() {
		n.Wait()
		close(metas)
	}()
	results := []MarkdownMeta{}
	for m := range metas {
		results = append(results, m)
	}
	return results
}

func doTraversal(dir string, n *sync.WaitGroup, metas chan<- MarkdownMeta) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		fname := filepath.Join(dir, entry.Name())
		switch {
		case strings.HasPrefix(entry.Name(), "."):
			// dot file, ignore
		case entry.IsDir():
			// dir, dive into
			n.Add(1)
			go doTraversal(fname, n, metas)
		case strings.HasSuffix(entry.Name(), ".md"):
			// a .md file, render it
			n.Add(1)
			go doRender(fname, n, metas)
		default:
			// static file, copy it
			copyFile(fname, filepath.Join(env.FrontEnd, fname[len(env.Config().ArticleRepo):]))
		}
	}
}

func doRender(fname string, n *sync.WaitGroup, metas chan<- MarkdownMeta) {
	defer n.Done()
	// acquire token
	renderSema <- struct{}{}
	defer func() {
		// release token
		<-renderSema
	}()
	m, err := NewMarkdown(fname)
	if err != nil {
		panic(fmt.Sprintf("render md %s fail, %v\n", fname, err))
	}
	b, err := m.Render()
	if err != nil {
		panic(fmt.Sprintf("render md %s fail, %v\n", fname, err))
	}
	// save the file as .html
	// TODO: transform it to ``index.html`` for simple use cases
	dest := filepath.Join(env.FrontEnd, m.Id+".html")
	ensureDir(dest)
	fmt.Printf("%s -> %s\n", fname, dest)
	err = ioutil.WriteFile(dest, b, 0644)
	if err != nil {
		panic(fmt.Sprintf("write file fail, %v\n", err))
	}
	metas <- m.MarkdownMeta
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
