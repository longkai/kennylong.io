package env

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestInitEnv(t *testing.T) {
	defer func() {
		if v := recover(); v != nil {
			t.Errorf("Init env fail, %v\n", v)
		}
	}()

	InitEnv("../testing_env.json")
	c := Config()

	if strings.HasSuffix(c.AccessToken, string(filepath.Separator)) {
		t.Errorf("repo path should not has Separator.\n")
	}
}

func TestIgnore(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{env.ArticleRepo + "/aaa", true},                              // none *.md
		{env.ArticleRepo + "/a.md", true},                             // a *.md
		{env.ArticleRepo + "/" + env.PublishDirs[0], false},           // sub dir
		{env.ArticleRepo + "/" + env.PublishDirs[0] + "/c.md", false}, // sub dir' s *.md
		{env.ArticleRepo, true},                                       // same
		{env.ArticleRepo + "/", true},                                 // same plus a slash
		{"a/b/c/.md", true},                                           // not in the repo
		{"", true},                                                    // empty
	}

	for _, test := range tests {
		if got := Ignored(test.input); got != test.want {
			t.Errorf("Ignored(%q) = %v", test.input, got)
		}
	}
}
