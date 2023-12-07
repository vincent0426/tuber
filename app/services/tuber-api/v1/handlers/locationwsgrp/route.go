package locationwsgrp

import (
	"net/http"

	"github.com/TSMC-Uber/server/business/core/locationws"
	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/foundation/logger"
	"github.com/TSMC-Uber/server/foundation/web"
)

type Config struct {
	Log  *logger.Logger
	Auth *auth.Auth
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	locationwsCore := locationws.NewCore(cfg.Log)

	// authen := mid.Authenticate(cfg.Auth)

	hdl := New(locationwsCore)
	app.Handle(http.MethodGet, version, "/ws/driver", hdl.DriverWebSocketHandler)
	app.Handle(http.MethodGet, version, "/ws/passenger", hdl.PassengerWebSocketHandler)
}
