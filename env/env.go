package env

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	FrontEnd = "../frontend" // where puts the static files
)

type Env struct {
	AccessToken string `json:"access_token"`
	ArticleRepo string `json:"article_repo"`
}

var (
	env *Env
)

func init() {
	// TODO: for simple testing, we hard code here
	InitEnv("../testing_env.json")
	// ensure front dir exist
	_, err := os.Stat(FrontEnd)
	if os.IsNotExist(err) {
		err = os.Mkdir(FrontEnd, 0755)
		if err != nil {
			panic(fmt.Sprintf("Mkdir %s fail, %v\n", FrontEnd))
		}
	}
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
}

func Config() Env {
	if env == nil {
		panic(fmt.Sprintf("plz call `InitEnv(string)` first"))
	}
	return *env
}
