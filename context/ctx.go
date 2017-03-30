package context

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Conf configuration
type Conf struct {
	Port        int      `yaml:"port"`
	RepoDir     string   `yaml:"repo_dir"`
	MediumToken string   `yaml:"medium_token"`
	GlobDocs    []string `yaml:"glob_docs"`
	SkipDirs    []string `yaml:"skip_dirs"`
	Github      struct {
		User        string `yaml:"user"`
		Repo        string `yaml:"repo"`
		HookSecret  string `yaml:"hook_secret"`
		AccessToken string `yaml:"access_token"`
	} `yaml:"github"`
	Meta struct {
		V             string
		B             string
		GA            string `json:"ga"`
		GF            bool   `json:"gf"`
		Origin        string `json:"origin"`
		Bio           string `json:"bio"`
		Link          string `json:"link"`
		Name          string `json:"name"`
		Title         string `json:"title"`
		Mail          string `json:"mail"`
		Github        string `json:"github"`
		Medium        string `json:"medium"`
		Twitter       string `json:"twitter"`
		Instagram     string `json:"instagram"`
		Stackoverflow string `json:"stackoverflow"`
	} `json:"meta"`
}

// Compile time variables.
var v, b string

func init() { fmt.Printf("Happy hacking :) Build ID: %s, Branch: %s\n", v, b) }

// NewConf new configuration properties for the given file path.
func NewConf(path string) (Conf, error) {
	var conf = Conf{}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, err
	}
	if err = yaml.Unmarshal(bytes, &conf); err != nil {
		return conf, err
	}
	conf.Meta.V, conf.Meta.B = v, b
	return conf, nil
}
