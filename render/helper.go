package render

import (
	"fmt"
	"time"
)

// DaysAgo days ago.
func DaysAgo(t time.Time) int { return int(time.Since(t).Hours() / 24) }

// Format simple datetime format.
func Format(t time.Time) string { return t.Format(time.RFC1123Z) }

// Tags _
func Tags(m []string) string {
	s := ""
	for i, v := range m {
		s += "#" + v
		if i != len(m)-1 {
			s += ", "
		}
	}
	if s[0] != '#' {
		// image
		return fmt.Sprintf("url('%s')", s)
	}
	return s
}

// BgImg s resolove is a color or a image attr
func BgImg(s string) bool {
	if len(s) == 0 { // `red` `green` is not allow, instead of hex
		return false // in case color write empty string which is ok
	}
	return s[0] != '#'
}
