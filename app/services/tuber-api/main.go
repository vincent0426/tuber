// @title TUber API
// @description This is the API for the TUber application.
// @version develop
// @host localhost:3000
// @BasePath /v1
// @schemes http https ws wss
package main

import (
	"os"

	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/cmd"
	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/cmd/all"
	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/cmd/chat"
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

	case "chat":
		if err := cmd.Main(build, chat.Routes()); err != nil {
			os.Exit(1)
		}

		// case "reporting":
		// 	if err := cmd.Main(build, reporting.Routes()); err != nil {
		// 		os.Exit(1)
		// 	}
	}
}
