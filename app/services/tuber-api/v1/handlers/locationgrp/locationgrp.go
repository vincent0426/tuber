// Package tripgrp maintains the group of handlers for trip access.
package locationgrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/TSMC-Uber/server/business/core/location"
	"github.com/TSMC-Uber/server/business/web/v1/paging"
	"github.com/TSMC-Uber/server/business/web/v1/response"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	location *location.Core
}

// New constructs a handlers for route access.
func New(location *location.Core) *Handlers {
	return &Handlers{
		location: location,
	}
}

// @Summary create a new location
// @Schemes
// @Description Create will add a location
// @Tags location
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param body body AppNewLocation true "New Location"
// @Success 201 {object} AppLocation "Location successfully created"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /locations [post]
func (h *Handlers) Create(ctx context.Context, c *gin.Context) error {
	var app AppNewLocation
	// Validate the request.
	if err := web.Decode(c, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	nc, err := toCoreNewLocation(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	location, err := h.location.Create(ctx, nc)
	if err != nil {
		return fmt.Errorf("create: location[%+v]: %w", location, err)
	}

	return web.Respond(ctx, c.Writer, toAppLocation(location), http.StatusCreated)
}

// @Summary get all locations
// @Schemes
// @Description Query will query locations
// @Tags location
// @Accept json
// @Produce json
// @Success 200 {object} AppLocation "Locations successfully queried"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /locations [get]
func (h *Handlers) Query(ctx context.Context, c *gin.Context) error {
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

	locations, err := h.location.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	items := make([]AppLocation, len(locations))
	for i, location := range locations {
		items[i] = toAppLocation(location)
	}

	total, err := h.location.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, c.Writer, paging.NewResponse(items, total, page.Number, page.RowsPerPage), http.StatusOK)
}

// @Summary get locations by id
// @Schemes
// @Description QueryByID will query locations by id
// @Tags location
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} AppLocation "query locations"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /locations/{id} [get]
func (h *Handlers) QueryByID(ctx context.Context, c *gin.Context) error {
	id := c.Param("id")

	qlocation, err := h.location.QueryByID(ctx, uuid.Must(uuid.Parse(id)))
	if err != nil {
		switch {
		case errors.Is(err, location.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: id[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, c.Writer, toAppLocation(qlocation), http.StatusOK)
}
