package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/longkai/xiaolongtongxue.com/config"
	"github.com/longkai/xiaolongtongxue.com/controller"
)

var (
	port = flag.Int(`port`, 1217, `HTTP listen port`)
	conf = flag.String(`conf`, `env.yaml`, `configuration file`)
)

func main() {
	flag.Parse()
	if err := config.Init(*conf); err != nil {
		log.Fatal(err)
	}
	controller.Ctrl()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
