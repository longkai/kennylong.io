package main

import (
	"fmt"
	"os"

	"github.com/longkai/xiaolongtongxue.com/helper"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s [src] [dest]\n", os.Args[0])
		os.Exit(1)
	}

	src, dest := os.Args[1], os.Args[2]
	if err := helper.JSON2YamlFile(src, dest); err != nil {
		panic(err)
	}

	fmt.Println("done.")
}
