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
	pt       PathTransformer
	skipDirs []string
	globDocs []string
}

// Scan a path to find documents.
func (s *DocScanner) Scan(path string) []string {
	var res []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		switch {
		case info.IsDir():
			for _, p := range s.skipDirs {
				// If match fail, don't skip, same below.
				if ok, _ := filepath.Match(filepath.Join(s.pt.baseDir, p), path); ok {
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
	return res
}
