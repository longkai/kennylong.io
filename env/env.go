package env

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Env struct {
	AccessToken string `json:"access_token"`
}

var (
	env *Env
)

func init() {
	// TODO: for simple testing, we hard code here
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
}

func Config() Env {
	if env == nil {
		panic(fmt.Sprintf("plz call `InitEnv(string)` first"))
	}
	return *env
}
