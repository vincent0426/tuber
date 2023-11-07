package authgrp

import (
	"net/http"

	aauth "github.com/TSMC-Uber/server/business/core/auth"
	"github.com/TSMC-Uber/server/business/core/auth/stores/authdb"
	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/business/web/v1/mid"
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

	authCore := aauth.NewCore(authdb.NewStore(cfg.Log, cfg.DB))

	authenGoogle := mid.AuthGoogle(cfg.Auth, authCore)

	hdl := New(authCore)
	app.Handle(http.MethodPost, version, "/auth/login", hdl.Login, authenGoogle)
	app.Handle(http.MethodPost, version, "/auth/logout", hdl.Logout, authenGoogle)
}
