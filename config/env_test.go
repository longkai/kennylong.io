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
	Env = &Configuration{Repo: "/a/b/c/", Ignored: []string{"ignore", "ignore.txt"}}
	t.Log(Env)

	tests := []struct {
		tag   string
		input string
		want  bool
	}{
		{"rootFile", "file.md", true},
		{"rootDir", "file", false},
		{"ignoredFile", "a/ignore.txt", true},
		{"ingoreDir", "x/y/z/ignore", true},
		{"normalFile", "x/y/z/balh.txt", false},
		{"normalDir", "x/y/z", false},
	}

	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			input := Env.Repo + test.input
			if got := IsIgnored(input); got != test.want {
				t.Errorf("IsIgnored(%q) = %v", input, got)
			}
		})
	}
}
