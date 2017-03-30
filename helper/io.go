package helper

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// ioSemas is a counting semaphore for limiting concurrency in io events.
var ioSemas = make(chan struct{}, 20)

// Dirents lists the entries of directory dir.
func Dirents(dir string) []os.FileInfo {
	ioSemas <- struct{}{}
	defer func() { <-ioSemas }()
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Print(err)
		return nil
	}
	return entries
}

// Cp from src to dest
func Cp(src, dest string) error {
	ioSemas <- struct{}{}
	defer func() { <-ioSemas }()
	// ensure parent existence
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	if err = r.Sync(); err != nil {
		return err
	}
	return nil
}

// Exists test a path existence on the OS.
func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return !os.IsNotExist(err)
	}
	return true
}
