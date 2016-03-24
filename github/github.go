package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/longkai/xiaolongtongxue.com/env"
	"github.com/longkai/xiaolongtongxue.com/render"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

func Hook(resp http.ResponseWriter, req *http.Request, requests chan<- render.Articles) {
	signature := req.Header.Get("X-Hub-Signature")
	delivery := req.Header.Get("X-GitHub-Delivery")
	event := req.Header.Get("X-GitHub-Event")
	// handle security problem
	defer req.Body.Close()
	if err := handleSecuriy(req.Body, signature); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("receive github webhook, event: %q, delivery: %q, signature: %q\n", event, delivery, signature)
	if event == "push" {
		// NOTE: we do a simple job, each time we receive a push hook, just pull the master brach, then render again
		go pull(requests)
	}
	// send pong message back to Github
	fmt.Fprint(resp, "thx :)")
}

func pull(requests chan<- render.Articles) {
	log.Println("executing shell command...")
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd %s; git pull;", env.Config().ArticleRepo))
	b, err := cmd.Output()
	if err != nil {
		log.Printf("git pull fail, %v\n", err)
		return
	}
	fmt.Printf("%s\n", b)
	// go render again :)
	go doRender(requests)
}

func doRender(requests chan<- render.Articles) {
	articles := render.Traversal(env.Config().ArticleRepo)
	log.Printf("hook -> pull -> reload %d articles :) \n", len(articles))
	requests <- articles
}

func handleSecuriy(reader io.Reader, signature string) error {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	if !checkMAC(b, signature) {
		return fmt.Errorf("signature checking fail")
	}
	return nil
}

func checkMAC(message []byte, messageMAC string) bool {
	mac := hmac.New(sha1.New, []byte(env.Config().HookSecret))
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return "sha1="+hex.EncodeToString(expectedMAC) == messageMAC
}
