package healthgrp

import (
	"context"
	"fmt"
	"net/http"

	_ "github.com/TSMC-Uber/server/app/services/tuber-api/v1/docs"
	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/foundation/logger"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log  *logger.Logger
	Auth *auth.Auth
	DB   *sqlx.DB
	Svc  string
}

// @Summary ping
// @Schemes
// @Description do ping
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {string} string "pong"
// @Router /ping [get]
func Pong(ctx context.Context, c *gin.Context) error {
	return web.Respond(ctx, c.Writer, "pong", http.StatusOK)
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

	app.Engine.GET("/v1/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.Handle(http.MethodGet, version, healthPath, Pong)
}
