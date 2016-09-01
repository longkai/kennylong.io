package config

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// Configuration configuration
type Configuration struct {
	Repo        string   `yaml:"repo"`
	Ignored     []string `yaml:"ignored"`
	HookSecret  string   `yaml:"hook_secret"`
	AccessToken string   `yaml:"access_token"`
}

var (
	// Env global environment
	Env *Configuration
)

// InitEnv _
func InitEnv(src string) error {
	bytes, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	Env = new(Configuration)
	if err = yaml.Unmarshal(bytes, Env); err != nil {
		return err
	}

	// adjuest for simply handling stuffs
	if !strings.HasSuffix(Env.Repo, "/") {
		Env.Repo += "/"
	}
	return nil
}

// IsIgnored the abs path if has the same prefix
func IsIgnored(path string) bool {
	path = path[len(Env.Repo):]
	// root dir *.md are igored by default
	if !strings.Contains(path, "/") && strings.HasSuffix(path, ".md") {
		return true
	}

	for _, v := range Env.Ignored {
		if strings.Contains(path, v) {
			return true
		}
	}

	return false
}
