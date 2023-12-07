package taskgrp

import (
	"net/http"

	"github.com/TSMC-Uber/server/business/core/task"
	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/foundation/logger"
	"github.com/TSMC-Uber/server/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log  *logger.Logger
	Auth *auth.Auth
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	// TODO: Authen who can access this API
	// authen := mid.Authenticate(cfg.Auth)

	taskCore := task.NewCore()
	hdl := New(taskCore)

	app.Handle(http.MethodPost, version, "/tasks", hdl.PublishSendEmailTask)

}
