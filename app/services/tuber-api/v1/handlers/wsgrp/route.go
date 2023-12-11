package wsgrp

import (
	"net/http"

	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/TSMC-Uber/server/business/core/user/stores/userdb"
	"github.com/TSMC-Uber/server/business/core/ws"
	"github.com/TSMC-Uber/server/business/core/ws/stores/wsdb"
	"github.com/TSMC-Uber/server/business/web/v1/auth"
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

	wsCore := ws.NewCore(wsdb.NewStore(cfg.Log, cfg.DB))
	usrCore := user.NewCore(userdb.NewStore(cfg.Log, cfg.DB))

	authen := mid.Authenticate(cfg.Auth)

	hdl := New(wsCore, usrCore)
	app.Handle(http.MethodGet, version, "/chat/ws", hdl.Connect, authen)
}
