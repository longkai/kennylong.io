package env

import "testing"

func TestInitEnv(t *testing.T) {
	saved := ensureFrontEndDir
	defer func() { ensureFrontEndDir = saved }()

	// stub
	ensureFrontEndDir = func(s string) error { return nil }

	src := "./testdata/env.yaml"

	if err := InitEnv(src); err != nil {
		t.Errorf("InitEnv(%q) = %v\n", err)
	}

	if env == nil {
		t.Errorf("InitEnv(%q), env = nil\n", src)
	}
}

func TestIgnore(t *testing.T) {
	saved := env
	defer func() { env = saved }()

	// stub
	env = &Env{AccessToken: "blah", ArticleRepo: "/a/b/c/", HookSecret: "sec", PublishDirs: []string{"1", "2"}}

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
