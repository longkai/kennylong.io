package config

import "testing"

func TestInit(t *testing.T) {
	saved1, saved2 := Env, regexps
	defer func() {
		Env = saved1
		regexps = saved2
	}()

	src := "./testdata/env.yml"

	if err := Init(src); err != nil {
		t.Errorf("Init(%q) = %v\n", src, err)
	}

	if Env == nil {
		t.Errorf("Init(%q), Env = nil\n", src)
	}

	// testify defauly ignore behavious
	tests := []struct {
		tag   string
		input string
		want  bool
	}{
		{"RootHidden", ".git", true},
		{"CustomizedHidden", ".md", true},
		{"NormalHidden", "a/c/.c", true},
		{"Normal", "a/b/c.c", false},
		{"Normal", "a/b/c.md", false},
	}

	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			input := Env.Repo + test.input
			if got := Ignored(input); got != test.want {
				t.Errorf("Ignored(%q) = %v", input, got)
			}
		})
	}
}

func TestIgnore(t *testing.T) {
	saved1, saved2 := Env, regexps
	defer func() {
		Env = saved1
		regexps = saved2
	}()

	// stub
	Env = &Configuration{Repo: "/a/b/c/", Ignores: []string{`/[^/]*\.c$`, `.*/ignore/.*`}}
	adjustEnv()

	tests := []struct {
		tag   string
		input string
		want  bool
	}{
		{"Normal", "path/to/sth", false},
		{"RootExt", "blah.c", true},
		{"RootHiddenFile", ".blah.c", true},
		{"RootHiddenDir", ".blah", true},
		{"HiddenFile", "x/y/.z", true},
		{"IgnoreDir", "ignore/a/b/c", true},
	}

	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			input := Env.Repo + test.input
			if got := Ignored(input); got != test.want {
				t.Errorf("Ignored(%q) = %v", input, got)
			}
		})
	}
}

func TestRoot(t *testing.T) {
	_env, _roots := Env, roots
	defer func() { Env, roots = _env, _roots }()

	// stub
	Env = &Configuration{Repo: "/a/b/c/", Ignores: []string{`/[^/]*\.c$`, `.*/ignore/.*`}}
	adjustEnv()
	roots = map[string]struct{}{"exist": struct{}{}}

	tests := []struct {
		tag, input, want string
	}{
		{`Normal`, `path/to/sth`, `path`},
		{`Exist`, `exist/to/sth`, ``},
		{`Ignored`, `ignore/to/sth`, ``},
	}

	for _, test := range tests {
		test := test
		t.Run(test.tag, func(t *testing.T) {
			path := Env.Repo + test.input
			if got := Root(path); got != test.want {
				t.Errorf("Root(%q) = %q, want %q", path, got, test.want)
			}
		})
	}
}
