package tripgrp

import (
	"net/http"

	"github.com/TSMC-Uber/server/business/core/trip"
	"github.com/TSMC-Uber/server/business/core/trip/stores/tripdb"
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

	tripCore := trip.NewCore(tripdb.NewStore(cfg.Log, cfg.DB))

	authen := mid.Authenticate(cfg.Auth)

	hdl := New(tripCore)
	app.Handle(http.MethodGet, version, "/trips", hdl.Query)
	app.Handle(http.MethodGet, version, "/trips/:id", hdl.QueryByID)
	app.Handle(http.MethodPost, version, "/trips", hdl.Create, authen)
	app.Handle(http.MethodPut, version, "/trips/:id", hdl.Update, authen)
	// app.Handle(http.MethodPost, version, "/trips/join", hdl.Join)
	// app.Handle(http.MethodDelete, version, "/users/:id", hdl.Delete)
	app.Handle(http.MethodGet, version, "/trips/my", hdl.QueryMyTrip, authen)
	app.Handle(http.MethodPost, version, "/trips/join", hdl.Join, authen)
}
