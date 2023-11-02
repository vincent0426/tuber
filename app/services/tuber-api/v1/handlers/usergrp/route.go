package usergrp

import (
	"net/http"

	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/TSMC-Uber/server/business/core/user/stores/userdb"
	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log  *zap.SugaredLogger
	Auth *auth.Auth
	DB   *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	// envCore := event.NewCore(cfg.Log)
	// usrCore := user.NewCore(cfg.Log, envCore, usercache.NewStore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB)))
	usrCore := user.NewCore(userdb.NewStore(cfg.Log, cfg.DB))

	hdl := New(usrCore)
	app.Handle(http.MethodGet, version, "/users", hdl.Query)
	app.Handle(http.MethodPost, version, "/users", hdl.Create)
	app.Handle(http.MethodPut, version, "/users/:id", hdl.Update)
	app.Handle(http.MethodDelete, version, "/users/:id", hdl.Delete)
}
