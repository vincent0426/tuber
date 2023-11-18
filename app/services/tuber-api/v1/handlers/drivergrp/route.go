package drivergrp

import (
	"net/http"

	"github.com/TSMC-Uber/server/business/core/driver"
	"github.com/TSMC-Uber/server/business/core/driver/stores/driverdb"
	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/business/web/v1/mid"
	"github.com/TSMC-Uber/server/foundation/logger"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

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

	// envCore := event.NewCore(cfg.Log)
	// usrCore := user.NewCore(cfg.Log, envCore, usercache.NewStore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB)))

	driverCore := driver.NewCore(driverdb.NewStore(cfg.Log, cfg.DB))

	authen := mid.Authenticate(cfg.Auth)

	hdl := New(driverCore)
	app.Handle(http.MethodPost, version, "/drivers", hdl.Create, authen)
	app.Handle(http.MethodGet, version, "/drivers", hdl.QueryAll)
	app.Handle(http.MethodGet, version, "/drivers/:id", hdl.QueryByID)
}
