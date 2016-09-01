package config

import "testing"

func TestInitEnv(t *testing.T) {
	src := "./testdata/env.yaml"

	if err := InitEnv(src); err != nil {
		t.Errorf("InitEnv(%q) = %v\n", err)
	}

	if Env == nil {
		t.Errorf("InitEnv(%q), Env = nil\n", src)
	}
}

func TestIgnore(t *testing.T) {
	saved := Env
	defer func() { Env = saved }()

	// stub
	Env = &Configuration{AccessToken: "blah", ArticleRepo: "/a/b/c/", HookSecret: "sec", PublishDirs: []string{"1", "2"}}

	tests := []struct {
		input string
		want  bool
	}{
		{Env.ArticleRepo + "/aaa", true},                              // none *.md
		{Env.ArticleRepo + "/a.md", true},                             // a *.md
		{Env.ArticleRepo + "/" + Env.PublishDirs[0], false},           // sub dir
		{Env.ArticleRepo + "/" + Env.PublishDirs[0] + "/c.md", false}, // sub dir' s *.md
		{Env.ArticleRepo, true},                                       // same
		{Env.ArticleRepo + "/", true},                                 // same plus a slash
		{"a/b/c/.md", true},                                           // not in the repo
		{"", true},                                                    // empty
	}

	for _, test := range tests {
		if got := Ignored(test.input); got != test.want {
			t.Errorf("Ignored(%q) = %v", test.input, got)
		}
	}
}
