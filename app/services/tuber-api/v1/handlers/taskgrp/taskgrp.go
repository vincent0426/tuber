package taskgrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TSMC-Uber/server/business/core/task"
	"github.com/TSMC-Uber/server/business/sys/mq"
	"github.com/TSMC-Uber/server/business/web/v1/response"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/gin-gonic/gin"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	task *task.Core
}

// New constructs a handlers for route access.
func New(task *task.Core) *Handlers {
	return &Handlers{
		task: task,
	}
}

func (h *Handlers) PublishSendEmailTask(ctx context.Context, c *gin.Context) error {
	var app AppPublishSendEmailTask

	if err := web.Decode(c, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	fmt.Println("------------------")
	fmt.Println(app.DelayTime)
	fmt.Println(app.Email)
	fmt.Println("------------------")
	mq.SendDelayMsg([]byte(app.Email), app.DelayTime)

	return web.Respond(ctx, c.Writer, app, http.StatusOK)
}
