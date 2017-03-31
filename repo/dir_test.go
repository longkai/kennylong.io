package repo_test

import (
	"fmt"
	"testing"

	"github.com/longkai/xiaolongtongxue.com/repo"
)

func TestAbs(t *testing.T) {
	var base = "/root"

	tests := []struct {
		input string
		want  string
	}{
		{"", base},
		{"a", fmt.Sprintf("%s/a", base)},
		{"a/", fmt.Sprintf("%s/a", base)},
		{"/a", fmt.Sprintf("%s/a", base)},
		{"/a/", fmt.Sprintf("%s/a", base)},
		{"a/b/c", fmt.Sprintf("%s/a/b/c", base)},
		{base, base},
		{"/root/a/b", "/root/a/b"},
	}
	dir := repo.Dir(base)
	for _, test := range tests {
		if got := dir.Abs(test.input); got != test.want {
			t.Errorf("dir.Abs(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestRel(t *testing.T) {
	var base = "/root"
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"/", ""},
		{"/tmp", ""},       // Don't allow outside base dir if absolute.
		{"a/b/c", "a/b/c"}, // Return same if it's relative already.
		{"/root", "."},
		{"/root/", "."},
		{"/root/a", "a"},
		{"/root/a/", "a"},
		{"/root/a/b", "a/b"},
	}
	dir := repo.Dir(base)

	for _, test := range tests {
		if got := dir.Rel(test.input); got != test.want {
			t.Errorf("dir.Rel(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestURLPath(t *testing.T) {
	var base = "/root"
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"/", ""},
		{"/outside", ""},
		{"/root", ""},           // URL path roots from base dir
		{"/root/a", ""},         // Exclude root.
		{"/root/a/b/c", "/a/b"}, // URL path trim the last segment.
		{"/root/a/b/c.ext", "/a/b"},
	}
	dir := repo.Dir(base)

	for _, test := range tests {
		if got := dir.URLPath(test.input); got != test.want {
			t.Errorf("dir.URLPath(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}
