package main

import (
	_ "myproject/boot"
	_ "myproject/router"

	"github.com/wanghaha-dev/gf/frame/g"
)

func main() {
	g.Server().Run()
}
