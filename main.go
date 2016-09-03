package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/longkai/xiaolongtongxue.com/config"
	"github.com/longkai/xiaolongtongxue.com/controller"
)

func main() {
	var env = `env.yaml` // def location
	if len(os.Args) > 1 {
		env = os.Args[1]
	}
	if err := config.Init(env); err != nil {
		log.Fatalf("config.Init(%q) fail: %v", env, err)
	}
	controller.Ctrl()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Env.Port), nil))
}
