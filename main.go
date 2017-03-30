/*

                            oooo$$$$$$$$$$$$oooo
                        oo$$$$$$$$$$$$$$$$$$$$$$$$o
                     oo$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$o         o$   $$ o$
     o $ oo        o$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$o       $$ $$ $$o$
  oo $ $ "$      o$$$$$$$$$    $$$$$$$$$$$$$    $$$$$$$$$o       $$$o$$o$
  "$$$$$$o$     o$$$$$$$$$      $$$$$$$$$$$      $$$$$$$$$$o    $$$$$$$$
    $$$$$$$    $$$$$$$$$$$      $$$$$$$$$$$      $$$$$$$$$$$$$$$$$$$$$$$
    $$$$$$$$$$$$$$$$$$$$$$$    $$$$$$$$$$$$$    $$$$$$$$$$$$$$  """$$$
     "$$$""""$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$     "$$$
      $$$   o$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$     "$$$o
     o$$"   $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$       $$$o
     $$$    $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$" "$$$$$$ooooo$$$$o
    o$$$oooo$$$$$  $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$   o$$$$$$$$$$$$$$$$$
    $$$$$$$$"$$$$   $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$     $$$$""""""""
   """"       $$$$    "$$$$$$$$$$$$$$$$$$$$$$$$$$$$"      o$$$
              "$$$o     """$$$$$$$$$$$$$$$$$$"$$"         $$$
                $$$o          "$$""$$$$$$""""           o$$$
                 $$$$o                                o$$$"
                  "$$$$o      o$$$$$$o"$$$$o        o$$$$
                    "$$$$$oo     ""$$$$o$$$$$o   o$$$$""
                       ""$$$$$oooo  "$$$o$$$$$$$$$"""
                          ""$$$$$$$oo $$$$$$$$$$
                                  """"$$$$$$$$$$$
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/longkai/xiaolongtongxue.com/context"
	"github.com/longkai/xiaolongtongxue.com/controller"
)

func main() {
	var path = "conf.yml" // Default in current dir.
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	conf, err := context.NewConf(path)
	if err != nil {
		log.Fatalf("config.NewConf(%q): %v", path, err)
	}
	controller.Ctrl(conf)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil))
}
