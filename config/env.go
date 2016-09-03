package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/longkai/xiaolongtongxue.com/helper"

	"regexp"

	"gopkg.in/yaml.v2"
)

// Configuration configuration
type Configuration struct {
	Port        int      `yaml:"port"`
	Repo        string   `yaml:"repo"`
	Ignores     []string `yaml:"ignores"` // regexp
	HookSecret  string   `yaml:"hook_secret"`
	AccessToken string   `yaml:"access_token"`
	Meta        struct {
		GA            string `json:"ga"`
		CDN           string `json:"cdn"`
		Domain        string `json:"domain"`
		Bio           string `json:"bio"`
		Link          string `json:"link"`
		Lang          string `json:"lang"`
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

var (
	// Env global environment
	Env *Configuration

	regexps []*regexp.Regexp
	mu      sync.Mutex // guards roots
	roots   map[string]struct{}
)

var adjustEnv = func() {
	// adjuest for simply handling path stuffs
	if !strings.HasSuffix(Env.Repo, "/") {
		Env.Repo += "/"
	}

	regexps = make([]*regexp.Regexp, 0, len(Env.Ignores))
	regexps = append(regexps, regexp.MustCompile(`/\.[^/]+$`))     // ignore hidden file/dir
	regexps = append(regexps, regexp.MustCompile(`/assets(/.*)?`)) // ignore assets dir
	for _, v := range Env.Ignores {
		regexps = append(regexps, regexp.MustCompile(v))
	}

	// store root dirs
	mu.Lock()
	defer mu.Unlock()
	for _, e := range helper.Dirents(Env.Repo) {
		name := e.Name()
		if !Ignored(filepath.Join(Env.Repo, name)) {
			roots[name] = struct{}{}
		}
	}
}

// Init configuration, must call it only once.
func Init(src string) error {
	bytes, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	Env = new(Configuration)
	if err = yaml.Unmarshal(bytes, Env); err != nil {
		return err
	}
	roots = make(map[string]struct{})
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

// Roots return repo root dirs.
func Roots() []string {
	mu.Lock()
	defer mu.Unlock()
	list := make([]string, 0, len(roots))
	for k := range roots {
		list = append(list, k)
	}
	return list
}

// Root testify the path is in a new root dir, otherwise return ""
func Root(path string) string {
	if Ignored(path) {
		return ""
	}

	path = path[len(Env.Repo):] // without leading `/`
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			path = path[:i]
			break
		}
	}

	mu.Lock()
	defer mu.Unlock()
	if _, ok := roots[path]; ok {
		return ""
	}

	roots[path] = struct{}{}
	return path
}
