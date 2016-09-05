package render

import (
	"strings"
	"testing"

	"github.com/longkai/xiaolongtongxue.com/config"
)

func TestParseID(t *testing.T) {
	saved := config.Env
	defer func() { config.Env = saved }()
	// stub
	base := "/path/to/repo"
	config.Env = &config.Configuration{Repo: base}
	tests := []struct {
		input, want string
	}{
		{base + "/a/b/c.md", "/a/b/"},
		{base + "/a", "/"},
		{base + "/", "/"},
		{"balabala/bala/", ""},
	}

	f := func(t *testing.T) {
		for _, test := range tests {
			if got := parseID(test.input); got != test.want {
				t.Errorf("parseID(%q) = %q, want %q\n", test.input, got, test.want)
			}
		}
	}

	t.Run("WithoutSuffix", f)
	config.Env.Repo += "/"
	t.Run("WithSuffix", f)
}

func TestAbsURLRegexp(t *testing.T) {
	tests := []struct {
		tag   string
		input string
		want  bool
	}{
		{`Normal`, `https://balabala.com/api`, true},
		{`NormalFile`, `https://balabala.com/p/index.html`, true},
		{`NormalIndex`, `https://balabala.com/`, true},
		{`NormalIndexWithoutSlash`, `https://balabala.com`, true},
		{`UpperHTTP`, `HTTP://BALA.com/abc`, true},
		{`NoprotocolWithPath`, `//balabala.com`, true},
		{`NoProtocol`, `//`, true}, // about:blank
		{`FileServer`, `///`, true},
		{`Root`, `/`, true},
		{`RootPath`, `/a/b/c`, true},
		{`Segment`, `#section`, false},
		{`Rel`, `a/b/c`, false},
	}

	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			if got := reAbsURL.MatchString(test.input); got != test.want {
				t.Errorf("reAbsURL.MatchString(%q) = %t", test.input, got)
			}
		})
	}
}

func TestLinkify(t *testing.T) {
	prefix := `/prefix/`
	var str = `
[link1](path/to/link)
[ link 2](path/to/link "Alt")
[你好](/)
`

	var want = `
[link1](/prefix/path/to/link)
[ link 2](/prefix/path/to/link "Alt")
[你好](/)
`

	got, err := linkify(strings.NewReader(str), []byte(prefix))
	if err != nil {
		t.Errorf("linkify(%s) fail: %v", str, err)
	}
	if string(got) != want {
		t.Errorf("linkify(%s) = %s, want %s", str, got, want)
	}
}
