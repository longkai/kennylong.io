package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Configuration configuration
type Configuration struct {
	HookSecret  string   `yaml:"hook_secret"`
	AccessToken string   `yaml:"access_token"`
	ArticleRepo string   `yaml:"article_repo"`
	PublishDirs []string `yaml:"publish_dirs"`
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
	err = yaml.Unmarshal(bytes, Env)
	if err != nil {
		return err
	}
	return nil
}

// Ignored file not in the PublishDirs or top level *.md file
func Ignored(path string) bool {
	rel, err := filepath.Rel(Env.ArticleRepo, path)
	if err != nil {
		return true
	}

	// ignore top level *.md file since it doesn't make any sense, we won't use it as /index.html
	if !strings.ContainsRune(rel, filepath.Separator) && strings.HasSuffix(path, ".md") {
		return true
	}

	for _, v := range Env.PublishDirs {
		if strings.HasPrefix(rel, v) {
			return false
		}
	}
	return true
}
