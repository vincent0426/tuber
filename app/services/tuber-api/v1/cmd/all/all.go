// Package all binds all the routes into the specified app.
package all

import (
	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/handlers/authgrp"
	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/handlers/healthgrp"
	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/handlers/locationgrp"
	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/handlers/tripgrp"
	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/handlers/usergrp"
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
	authgrp.Routes(app, authgrp.Config{
		Log:  cfg.Log,
		Auth: cfg.Auth,
		DB:   cfg.DB,
	})
	usergrp.Routes(app, usergrp.Config{
		Log:  cfg.Log,
		Auth: cfg.Auth,
		DB:   cfg.DB,
	})
	tripgrp.Routes(app, tripgrp.Config{
		Log:  cfg.Log,
		Auth: cfg.Auth,
		DB:   cfg.DB,
	})
	locationgrp.Routes(app, locationgrp.Config{
		Log:  cfg.Log,
		Auth: cfg.Auth,
		DB:   cfg.DB,
	})
}
