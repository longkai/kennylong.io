package render

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestMarkdownRederRead(t *testing.T) {
	m := New("hello, world")
	if _, err := ioutil.ReadAll(m); err != nil {
		t.Errorf("reading fail, %v\n", err)
	}
}

func TestMarkdownRender(t *testing.T) {
	m := New("Hello world github/longkai#1 **cool**, and #1!")
	if _, err := m.Render(); err != nil {
		t.Errorf("render fail, %v\n", err)
	}
}

func TestMetadataRegexp(t *testing.T) {
	s := fmt.Sprintf("### EOF \n%sjson\n%s\n%s", CODE_BLOCK, `{
	"key": "value"
}`, CODE_BLOCK)
	if !metaRegexp.MatchString(s) {
		t.Errorf("%s not match json block!\n", s)
	}
	ss := metaRegexp.FindStringSubmatch(s)
	fmt.Println(ss[1])
}

func TestTitleRegexp(t *testing.T) {
	s := "Hello, World \n====\n"
	if !titleRegexp.MatchString(s) {
		t.Errorf("%s not match title block!\n", s)
	}
}

func TestParseTitle(t *testing.T) {
	want := "Hello, Wrold"
	s := want + " \n===\n\n"
	fmt.Println(s)
	title, b := parseTitle([]byte(s))
	if title != want {
		t.Errorf("want %s, but got %s\n", want, title)
	}

	if len(b) != 0 {
		t.Errorf("result should be 0, got %d \n", len(b))
	}
}

func TestSeparateMetaAndText(t *testing.T) {
	s := fmt.Sprintf("heheh### EOF \n%sjson\n%s\n%s", CODE_BLOCK, `{
		"date": "2006-01-02T15:04:05+07:00"
}`, CODE_BLOCK)
	text, _ := separateTextAndMeta([]byte(s))
	if strings.Contains(text, "EOF") {
		t.Errorf("should not contain **EOF**, got %s\n", text)
	}
	fmt.Println(text)
}

func TestTrimExt(t *testing.T) {
	cases := []string{
		"a.md",
		"a.c.md",
		".aa.md",
		"abc.",
		"abc",
		".",
	}
	wanted := []string{
		"a",
		"a.c",
		".aa",
		"abc",
		"abc",
		"",
	}
	for i, s := range cases {
		if v := trimExt(cases[i]); v != wanted[i] {
			t.Errorf("case %s, want %s, got %s\n", s, wanted[i], v)
		}
	}
}

/*
func TestNewMarkDown(t *testing.T) {
	// TODO: just a simple testing... any better way to do?
	m, _ := NewMarkdown("path/to/md")
	fmt.Println(m)
}
*/
