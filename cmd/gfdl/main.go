package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	re = regexp.MustCompile(`https?://.+\.\w+`)
	wg = sync.WaitGroup{}
)

func main() {
	if len(os.Args) != 3 {
		handle(fmt.Errorf("usage: %s urlsrc dest", os.Args[0]))
	}
	gfdl(os.Args[1], os.Args[2])
}

func handle(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func gfdl(src, dst string) {
	dir := filepath.Dir(dst)
	if err := mkdirs(dir); err != nil {
		handle(err)
	}
	b, err := get(src)
	if err != nil {
		handle(err)
	}
	ch := make(chan string)
	for _, item := range re.FindAll(b, -1) {
		wg.Add(1)
		url := string(item)
		go func() {
			defer wg.Done()
			bytes, err := get(url)
			if err != nil {
				handle(err)
			}
			if err := ioutil.WriteFile(filepath.Join(dir, url[strings.LastIndexByte(url, '/')+1:]), bytes, 0644); err != nil {
				handle(err)
			}
			ch <- url
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for url := range ch {
		b = bytes.Replace(b, []byte(url), []byte(url[strings.LastIndexByte(url, '/')+1:]), -1)
	}
	if err := ioutil.WriteFile(dst, b, 0644); err != nil {
		handle(err)
	}
	fmt.Println("done :)")
}

func mkdirs(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return os.MkdirAll(dir, 0755)
	}
	if !info.IsDir() {
		return fmt.Errorf("parent %q is not directory", dir)
	}
	return nil
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %q status: %d", url, resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
