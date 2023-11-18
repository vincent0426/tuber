// Package tripgrp maintains the group of handlers for trip access.
package drivergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/TSMC-Uber/server/business/core/driver"
	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/business/web/v1/paging"
	"github.com/TSMC-Uber/server/business/web/v1/response"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/gin-gonic/gin"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	driver *driver.Core
}

// New constructs a handlers for route access.
func New(driver *driver.Core) *Handlers {
	return &Handlers{
		driver: driver,
	}
}

// Create adds a new driver to the system.
func (h *Handlers) Create(ctx context.Context, c *gin.Context) error {
	userID := auth.GetUserID(ctx)
	var app AppNewDriver
	// Validate the request.
	if err := web.Decode(c, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	nc, err := toCoreNewDriver(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	nc.UserID = userID

	driver, err := h.driver.Create(ctx, nc)
	if err != nil {
		return fmt.Errorf("create: driver[%+v]: %w", driver, err)
	}

	return web.Respond(ctx, c.Writer, toAppDriver(driver), http.StatusCreated)
}

func (h *Handlers) QueryAll(ctx context.Context, c *gin.Context) error {
	page, err := paging.ParseRequest(c.Request)
	if err != nil {
		return err
	}

	filter, err := parseFilter(c.Request)
	if err != nil {
		return err
	}

	orderBy, err := parseOrder(c.Request)
	if err != nil {
		return err
	}

	trips, err := h.driver.QueryAll(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	items := make([]AppDriver, len(trips))
	for i, trip := range trips {
		items[i] = toAppDriver(trip)
	}

	total, err := h.driver.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, c.Writer, paging.NewResponse(items, total, page.Number, page.RowsPerPage), http.StatusOK)
}

func (h *Handlers) QueryByID(ctx context.Context, c *gin.Context) error {
	id := c.Param("id")

	qdriver, err := h.driver.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, driver.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: id[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, c.Writer, toAppDriver(qdriver), http.StatusOK)
}
