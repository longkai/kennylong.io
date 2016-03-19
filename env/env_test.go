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
