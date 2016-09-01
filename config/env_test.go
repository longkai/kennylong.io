package config

import "testing"

func TestInitEnv(t *testing.T) {
	saved1, saved2 := Env, regexps
	defer func() {
		Env = saved1
		regexps = saved2
	}()

	src := "./testdata/env.yaml"

	if err := InitEnv(src); err != nil {
		t.Errorf("InitEnv(%q) = %v\n", err)
	}

	if Env == nil {
		t.Errorf("InitEnv(%q), Env = nil\n", src)
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
