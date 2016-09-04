package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/longkai/xiaolongtongxue.com/config"
	"github.com/longkai/xiaolongtongxue.com/controller"
)

var rev string

func main() {
	fmt.Printf("Happy hacking:) Build ID: %q\n", rev)
	var env = `env.yaml` // def location
	if len(os.Args) > 1 {
		env = os.Args[1]
	}
	if err := config.Init(env, rev); err != nil {
		log.Fatalf("config.Init(%q) fail: %v", env, err)
	}
	controller.Ctrl()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Env.Port), nil))
}
