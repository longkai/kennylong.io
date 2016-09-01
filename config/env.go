package config

import (
	"io/ioutil"
	"strings"

	"regexp"

	"gopkg.in/yaml.v2"
)

// Configuration configuration
type Configuration struct {
	Repo        string   `yaml:"repo"`
	Ignores     []string `yaml:"ignores"` // regexp
	HookSecret  string   `yaml:"hook_secret"`
	AccessToken string   `yaml:"access_token"`
}

var (
	// Env global environment
	Env     *Configuration
	regexps []*regexp.Regexp
)

var adjustEnv = func() {
	// adjuest for simply handling path stuffs
	if !strings.HasSuffix(Env.Repo, "/") {
		Env.Repo += "/"
	}

	regexps = make([]*regexp.Regexp, 0, len(Env.Ignores))
	regexps = append(regexps, regexp.MustCompile(`/\.[^/]+$`)) // ignore hidden file/dir
	for _, v := range Env.Ignores {
		regexps = append(regexps, regexp.MustCompile(v))
	}
}

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
	adjustEnv()
	return nil
}

// Ignored the regexp matches
func Ignored(path string) bool {
	path = path[len(Env.Repo)-1:] // keep leading `/`
	for _, re := range regexps {
		if re.MatchString(path) {
			return true
		}
	}
	return false
}
