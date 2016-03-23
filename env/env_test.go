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
	cases := []string{
		env.ArticleRepo + "/aaa",                             // none *.md
		env.ArticleRepo + "/a.md",                            // a *.md
		env.ArticleRepo + "/" + env.PublishDirs[0],           // sub dir
		env.ArticleRepo + "/" + env.PublishDirs[0] + "/c.md", // sub dir' s *.md
		env.ArticleRepo,                                      // same
		env.ArticleRepo + "/",                                // same plus a slash
		"a/b/c/.md",                                          // not in the repo
		"",                                                   // empty
	}

	wanted := []bool{
		true,
		true,
		false,
		false,
		true,
		true,
		true,
		true,
	}

	for i, s := range cases {
		if v := Ignored(cases[i]); v != wanted[i] {
			t.Errorf("case %s, want %t, got %t\n", s, wanted[i], v)
		}
	}
}
