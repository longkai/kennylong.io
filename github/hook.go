package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/longkai/xiaolongtongxue.com/git"
)

// Callback after pull and diff success.
type Callback func(adds, mods, dels []string)

var (
	repo     string
	token    string
	secret   string
	callback Callback
)

// Init Github service
func Init(hookURL string, _repo string, _secret string, _token string, cb Callback) {
	repo, secret, token, callback = _repo, _secret, _token, cb
	http.HandleFunc(hookURL, hook)
}

// Hook github webhook service.
func hook(w http.ResponseWriter, r *http.Request) {
	event := r.Header.Get("X-GitHub-Event")
	signature := r.Header.Get("X-Hub-Signature")
	delivery := r.Header.Get("X-GitHub-Delivery")
	// handle security
	defer r.Body.Close()
	if err := handleSecurity(r.Body, signature); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	log.Printf("receive github webhook, event: %q, delivery: %q, signature: %q", event, delivery, signature)

	if event == "push" {
		go handleHook()
	}
	// send pong message back to Github
	fmt.Fprint(w, "thx :)")
}

var handleHook = func() {
	// 1. get current HEAD hash
	v, err := git.Rev(repo)
	if err != nil {
		log.Printf("git.Rev(%q) fail: %v", repo, err)
		return
	}
	// 2. pull the latest content
	if err = git.Pull(repo); err != nil {
		log.Printf("`git pull` fail: %v", err)
		return
	}
	// 3. diff changes
	a, m, d, err := git.Diff(repo, v)
	if err != nil {
		log.Printf("`git diff` fail: %v", err)
		return
	}
	// 4. let sb. known what changed
	callback(a, m, d)
}

var handleSecurity = func(reader io.Reader, signature string) error {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	if !checkMAC(b, signature) {
		return fmt.Errorf("signature checking fail")
	}
	return nil
}

var checkMAC = func(message []byte, messageMAC string) bool {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return "sha1="+hex.EncodeToString(expectedMAC) == messageMAC
}
