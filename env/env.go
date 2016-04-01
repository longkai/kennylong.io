package env

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	GEN      = "gen" // where puts the gen articles files, in the root of this repo.
	Template = "templ"
)

type Env struct {
	HookSecret  string   `json:"hook_secret"`
	AccessToken string   `json:"access_token"`
	ArticleRepo string   `json:"article_repo"`
	PublishDirs []string `json:"publish_dirs"`
}

var (
	env *Env
)

func init() {
	// TODO: for simple testing, we hard code here
	defer func() {
		if v := recover(); v != nil {
			fmt.Fprintf(os.Stderr, "**env init() fail, it is okay if you are NOT in testing environment.**\n")
		}
	}()
	InitEnv("../testing_env.json")
}

func InitEnv(src string) {
	bytes, err := ioutil.ReadFile(src)
	if err != nil {
		panic(fmt.Sprintf("Init env from src %q fail, %v\n", src, err))
	}
	err = json.Unmarshal(bytes, &env)
	if err != nil {
		panic(fmt.Sprintf("Unmarshal env json fail, %v\n", err))
	}

	// trim the last `/`
	if env.ArticleRepo[len(env.ArticleRepo)-1] == filepath.Separator {
		env.ArticleRepo = env.ArticleRepo[:len(env.ArticleRepo)-1]
	}

	// ensure front dir exist
	_, err = os.Stat(GEN)
	if os.IsNotExist(err) {
		err = os.Mkdir(GEN, 0755)
		if err != nil {
			panic(fmt.Sprintf("Mkdir %s fail, %v\n", GEN))
		}
	}
}

func Config() Env {
	if env == nil {
		panic(fmt.Sprintf("plz call `InitEnv(string)` first"))
	}
	return *env
}

// Ignore file not in the PublishDirs or top level *.md file
func Ignored(path string) bool {
	rel, err := filepath.Rel(env.ArticleRepo, path)
	if err != nil {
		log.Println(err)
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
