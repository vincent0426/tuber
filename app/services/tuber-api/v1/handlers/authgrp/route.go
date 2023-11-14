package authgrp

import (
	"net/http"

	aauth "github.com/TSMC-Uber/server/business/core/auth"
	"github.com/TSMC-Uber/server/business/core/auth/stores/authdb"
	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/TSMC-Uber/server/business/core/user/stores/userdb"
	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/business/web/v1/mid"
	"github.com/TSMC-Uber/server/foundation/logger"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log     *logger.Logger
	Auth    *auth.Auth
	DB      *sqlx.DB
	RedisDB struct {
		Master  *redis.Client
		Replica *redis.Client
	}
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	userCore := user.NewCore(userdb.NewStore(cfg.Log, cfg.DB))

	authCore := aauth.NewCore(authdb.NewStore(cfg.Log, cfg.DB))

	authenGoogle := mid.AuthGoogle(cfg.Auth, authCore)

	hdl := New(authCore, userCore)
	app.Handle(http.MethodPost, version, "/auth/login", hdl.Login, authenGoogle)
	app.Handle(http.MethodPost, version, "/auth/logout", hdl.Logout, authenGoogle)
}
