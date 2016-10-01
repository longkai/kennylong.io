package controller

import (
	"bytes"
	"html/template"
	"time"
)

// DaysAgo days ago.
func DaysAgo(t time.Time) int { return int(time.Since(t).Hours() / 24) }

// Format simple datetime format.
func Format(t time.Time) string { return t.Format(time.RFC1123Z) }

// Tags #tag1, #tag2, ... #tagn
func Tags(m []string) string {
	buf := new(bytes.Buffer)
	for _, v := range m {
		buf.WriteByte('#')
		buf.WriteString(v)
		buf.WriteString(`, `)
	}
	buf.Truncate(buf.Len() - 2) // drop last `, `
	return buf.String()
}

// TransformCDN to cdn href
func TransformCDN(href string) template.URL {
	if cdn == "" {
		return template.URL(href)
	}
	return template.URL(cdn + revAsset(href))
}
