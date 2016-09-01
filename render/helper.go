package render

import (
	"strings"
	"time"
)

// DaysAgo days ago.
func DaysAgo(t time.Time) int { return int(time.Since(t).Hours() / 24) }

// Format simple datetime format.
func Format(t time.Time) string { return t.Format(time.RFC1123Z) }

// HasColor _
func HasColor(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] == '#'
}

// HasImage _
func HasImage(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] != '#'
}

// Tags _
func Tags(m []string) string {
	s := ""
	for i, v := range m {
		s += "#" + v
		if i != len(m)-1 {
			s += ", "
		}
	}
	return s
}

// IsRelImage _
func IsRelImage(s string) bool {
	if len(s) == 0 {
		return false
	}
	if s[0] == '/' || strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return false
	}
	return true
}
