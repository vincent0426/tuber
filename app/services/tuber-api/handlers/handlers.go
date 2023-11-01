package handlers

import (
	"net/http"
	"os"

	health "github.com/TSMC-Uber/server/app/services/tuber-api/handlers/v1/healthgrp"
	"github.com/TSMC-Uber/server/app/services/tuber-api/handlers/v1/usergrp"
	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/TSMC-Uber/server/business/core/user/stores/userdb"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	// Auth     *auth.Auth
	DB *sqlx.DB
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(cfg.Shutdown)

	health := health.New()
	app.GET("/health", health.Status)

	usrCore := user.NewCore(userdb.NewStore(cfg.Log, cfg.DB))

	ugh := usergrp.New(usrCore)

	app.Handle(http.MethodGet, "/users", ugh.Query)

	return app
}
