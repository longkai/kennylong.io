package repo

import (
	"path/filepath"
)

// PathTransformer transform path from abs to rel and vice-versa.
type PathTransformer struct {
	baseDir string
}

// Rel return the relative path for the given path, if not relative return "".
func (p PathTransformer) Rel(path string) string {
	if filepath.IsAbs(path) {
		rel, err := filepath.Rel(p.baseDir, path)
		if err != nil {
			return ""
		}
		return rel
	}
	return path
}

// Abs return the absolute path for the given path, if path is absolute but not in baseDir return "".
func (p PathTransformer) Abs(path string) string {
	if filepath.IsAbs(path) {
		if _, err := filepath.Rel(p.baseDir, path); err != nil {
			return ""
		}
		return path
	}
	return filepath.Join(p.baseDir, path)
}

// URLPath return the HTTP URL path for the file system path.
func (p PathTransformer) URLPath(path string) string {
	rel := p.Rel(path)
	if rel == "" {
		return ""
	}
	return filepath.Join("/", filepath.Dir(rel))
}

// NewPathTransformer create a path transformer for the given base dir.
func NewPathTransformer(baseDir string) PathTransformer {
	return PathTransformer{baseDir: baseDir}
}
