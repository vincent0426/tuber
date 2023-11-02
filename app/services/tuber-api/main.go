package main

import (
	"os"

	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/cmd"
	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/cmd/all"
)

/*
	Need to figure out timeouts for http service.
*/

var build = "develop"
var routes = "all" // go build -ldflags "-X main.routes=crud"

func main() {
	switch routes {
	case "all":
		if err := cmd.Main(build, all.Routes()); err != nil {
			os.Exit(1)
		}

		// case "crud":
		// 	if err := cmd.Main(build, crud.Routes()); err != nil {
		// 		os.Exit(1)
		// 	}

		// case "reporting":
		// 	if err := cmd.Main(build, reporting.Routes()); err != nil {
		// 		os.Exit(1)
		// 	}
	}
}
