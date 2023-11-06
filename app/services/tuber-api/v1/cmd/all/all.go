// Package all binds all the routes into the specified app.
package all

import (
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
	usergrp.Routes(app, usergrp.Config{
		Log:  cfg.Log,
		Auth: cfg.Auth,
		DB:   cfg.DB,
	})
}
