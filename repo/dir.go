package repo

import (
	"path/filepath"
	"strings"
)

// Dir transforms path from absolute to relative and vice-versa,
// from the repository base dir perspective.
type Dir string

// Rel returns the relative path for the given path,
// return same if path is already relative,
// if outside of the base dir return "".
func (p Dir) Rel(path string) string {
	if filepath.IsAbs(path) {
		rel, err := filepath.Rel(string(p), path)
		if err != nil {
			return ""
		}
		// Reject dir outside of base dir.
		if strings.HasPrefix(rel, "..") {
			return ""
		}
		return rel
	}
	return path
}

// Abs returns the absolute path in the prepend the base dir.
// Returns "" if the given path is absolute but not on the base dir.
func (p Dir) Abs(path string) string {
	if filepath.IsAbs(path) {
		if strings.HasPrefix(filepath.Clean(path), string(p)) {
			return path
		}
	}
	return filepath.Join(string(p), path)
}

// URLPath returns the HTTP URL path at the application level.
// Basically, return the dir of the path relative to base dir.
// Any path outside of base dir or equal to base dir will return "".
func (p Dir) URLPath(path string) string {
	rel := p.Rel(path)
	// Reject outside base dir and base dir.
	if rel == "" || rel == "." {
		return ""
	}
	// Reject top level dir, e.g., $ROOT/README.md will be reject,
	// $ROOT/dir/README.md will not, however.
	dir := filepath.Dir(rel)
	if dir == "." {
		return ""
	}
	return filepath.Join("/", dir)
}
