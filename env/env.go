package env

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	// Gen where puts the gen articles files, in the root of this repo.
	Gen = "_gen"
	// Template where the template located.
	Template = "templ"
)

// Env configuration
type Env struct {
	HookSecret  string   `yaml:"hook_secret"`
	AccessToken string   `yaml:"access_token"`
	ArticleRepo string   `yaml:"article_repo"`
	PublishDirs []string `yaml:"publish_dirs"`
}

var (
	env *Env
)

var ensureFrontEndDir = func(dirName string) error {
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		return os.Mkdir(dirName, 0755)
	}
	return nil
}

// InitEnv _
func InitEnv(src string) error {
	bytes, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, &env)
	if err != nil {
		return err
	}

	// ensure front-end dir exist
	return ensureFrontEndDir(Gen)
}

// Config get config
func Config() Env {
	if env == nil {
		panic(fmt.Sprintf("plz call `InitEnv(string)` first"))
	}
	return *env
}

// Ignored file not in the PublishDirs or top level *.md file
func Ignored(path string) bool {
	rel, err := filepath.Rel(env.ArticleRepo, path)
	if err != nil {
		return true
	}

	// ignore top level *.md file since it doesn't make any sense, we won't use it as /index.html
	if !strings.ContainsRune(rel, filepath.Separator) && strings.HasSuffix(path, ".md") {
		return true
	}

	for _, v := range env.PublishDirs {
		if strings.HasPrefix(rel, v) {
			return false
		}
	}
	return true
}
