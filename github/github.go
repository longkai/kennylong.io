package github

import (
  "bytes"
  "fmt"
  "log"
  "net/http"
  "os/exec"
  "github.com/longkai/xiaolongtongxue.com/render"
  "github.com/longkai/xiaolongtongxue.com/env"
)

func Hook(resp http.ResponseWriter, req *http.Request, requests chan <- render.Articles) {
  event := req.Header.Get("X-GitHub-Event")
  delivery := req.Header.Get("X-GitHub-Delivery")
  signature := req.Header.Get("X-Hub-Signature")
  log.Printf("receive github webhook, event: %q, delivery: %q, signature: %q\n", event, delivery, signature)
  if event == "push" {
    // NOTE: we do a simple job, each time we receive a push hook, just pull the master brach, then render again
    go pull(requests)
  }
  // send pong message back to Github
  fmt.Fprint(resp, "thx :)")
}

func pull(requests chan <- render.Articles) {
  cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd %s; git pull;", env.Config().ArticleRepo))
  b, err := cmd.Output()
  if err != nil {
    log.Printf("git pull fail, %v\n", err)
  }
  if index := bytes.LastIndex([]byte("Already up-to-date"), b); index != -1 {
    // go render again :)
    go doRender(requests)
    fmt.Printf("%s", b)
  } else {
    fmt.Println("alreay up-to-date. no need to fetch again...")
  }
}

func doRender(requests chan <- render.Articles) {
  articles := render.Traversal(env.Config().ArticleRepo)
  log.Printf("hook -> pull -> reload %d articles :) \n", len(articles))
  requests <- articles
}
