package healthgrp

import (
	"context"
	"net/http"

	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/gin-gonic/gin"
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

	app.Handle(http.MethodGet, version, "/ping", func(ctx context.Context, c *gin.Context) error {
		return web.Respond(ctx, c.Writer, "pong", http.StatusOK)
	})
}
