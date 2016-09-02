package render

import (
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
