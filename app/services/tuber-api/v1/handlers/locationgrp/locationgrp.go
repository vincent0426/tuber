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

// Create adds a new trip to the system.
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

// QueryAll returns a list of users with paging.
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

	locations, err := h.location.QueryAll(ctx, filter, orderBy, page.Number, page.RowsPerPage)
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

// QueryByID returns a trip by its ID.
func (h *Handlers) QueryByID(ctx context.Context, c *gin.Context) error {
	id := c.Param("id")

	qlocation, err := h.location.QueryByID(ctx, id)
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

// func (h *Handlers) Join(ctx context.Context, c *gin.Context) error {
// 	var app AppJoinTrip
// 	// Validate the request.
// 	if err := web.Decode(c, &app); err != nil {
// 		return response.NewError(err, http.StatusBadRequest)
// 	}

// 	nc, err := toCoreJoinTrip(app)
// 	if err != nil {
// 		return response.NewError(err, http.StatusBadRequest)
// 	}

// 	trip, err := h.trip.Join(ctx, nc)
// 	if err != nil {
// 		if errors.Is(err, user.ErrUniqueEmail) {
// 			return response.NewError(err, http.StatusConflict)
// 		}
// 		return fmt.Errorf("join: usr[%+v]: %w", trip, err)
// 	}

// 	return web.Respond(ctx, c.Writer, toAppTrip(trip), http.StatusCreated)
// }
