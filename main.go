package main

import (
	"flag"

	"github.com/Zioyi/temi/pkg/http"
	"github.com/Zioyi/temi/pkg/terminal"
)

var debugUrl = flag.String("url", "", "url for debug/vars, e.g. http://localhost/debug/vars")

func main() {
	flag.Parse()
	if len(*debugUrl) == 0 {
		flag.Usage()
		return
	}

	terminal.Run(http.NewMemStatsLoader(*debugUrl))
}
