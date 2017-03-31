package repo

import (
	"log"
	"os"
	"path/filepath"
)

// Scanner a location for documents.
type Scanner interface {
	Scan(path string) []string
}

// DocScanner scan a path for documents.
type DocScanner struct {
	dir      Dir
	skipDirs []string
	globDocs []string
}

// Scan a path to find documents.
func (s *DocScanner) Scan(path string) []string {
	var res []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		switch {
		case info == nil:
			log.Printf("info of %q is nil, err: %v", path, err)
			// Nope.
		case info.IsDir():
			for _, p := range s.skipDirs {
				// If match fail, don't skip, same below.
				if ok, _ := filepath.Match(filepath.Join(string(s.dir), p), path); ok {
					log.Printf("skip scan path: %q", path)
					return filepath.SkipDir
				}
			}
		default:
			for _, p := range s.globDocs {
				if ok, _ := filepath.Match(p, info.Name()); ok {
					res = append(res, path)
					break
				}
			}
		}
		return nil
	})
	log.Printf("scan %s, got %d", path, len(res))
	return res
}
