package helper

import (
	"io/ioutil"
	"log"
	"os"
)

// dirSema is a counting semaphore for limiting concurrency in dirents.
var dirSema = make(chan struct{}, 20)

// Dirents lists the entries of directory dir.
func Dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Print(err)
		return nil
	}
	return entries
}
