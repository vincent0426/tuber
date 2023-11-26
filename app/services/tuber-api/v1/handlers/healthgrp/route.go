package healthgrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/foundation/logger"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log  *logger.Logger
	Auth *auth.Auth
	DB   *sqlx.DB
	Svc  string
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	var healthPath string
	if cfg.Svc == "" {
		healthPath = "/ping"
	} else {
		healthPath = fmt.Sprintf("/%s/ping", cfg.Svc)
	}
	app.Handle(http.MethodGet, version, healthPath, func(ctx context.Context, c *gin.Context) error {
		return web.Respond(ctx, c.Writer, "pong", http.StatusOK)
	})
}
