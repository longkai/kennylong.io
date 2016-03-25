package render

import (
	"fmt"
	"strings"
	"testing"
)

func TestMetadataRegexp(t *testing.T) {
	s := fmt.Sprintf("### EOF \n%sjson\n%s\n%s", fenced_block, `{
	"key": "value"
}`, fenced_block)
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
	s := fmt.Sprintf("heheh### EOF \n%sjson\n%s\n%s", fenced_block, `{
		"date": "2006-01-02T15:04:05+07:00"
}`, fenced_block)
	text, _ := separateTextAndMeta([]byte(s))
	if strings.Contains(text, "EOF") {
		t.Errorf("should not contain **EOF**, got %s\n", text)
	}
	fmt.Println(text)
}

func TestTrimBasename(t *testing.T) {
	cases := []string{
		"a/a.md",
		"b/a.c.md",
		"abc.md",
		".",
	}
	wanted := []string{
		"a",
		"b",
		"abc.md",
		".",
	}
	for i, s := range cases {
		if v := trimBasename(cases[i]); v != wanted[i] {
			t.Errorf("case %s, want %s, got %s\n", s, wanted[i], v)
		}
	}
}
