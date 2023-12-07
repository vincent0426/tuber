package locationserver

import (
	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/handlers/healthgrp"
	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/handlers/locationwsgrp"
	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/handlers/taskgrp"
	v1 "github.com/TSMC-Uber/server/business/web/v1"
	"github.com/TSMC-Uber/server/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouterAdder interface.
func (add) Add(app *web.App, cfg v1.APIMuxConfig) {
	healthgrp.Routes(app, healthgrp.Config{
		Log: cfg.Log,
	})
	taskgrp.Routes(app, taskgrp.Config{
		Log:  cfg.Log,
		Auth: cfg.Auth,
	})
	locationwsgrp.Routes(app, locationwsgrp.Config{
		Log:  cfg.Log,
		Auth: cfg.Auth,
	})
}
