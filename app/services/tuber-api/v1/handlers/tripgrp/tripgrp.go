// Package tripgrp maintains the group of handlers for trip access.
package tripgrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/TSMC-Uber/server/business/core/trip"
	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/business/web/v1/paging"
	"github.com/TSMC-Uber/server/business/web/v1/response"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/gin-gonic/gin"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	trip *trip.Core
}

// New constructs a handlers for route access.
func New(trip *trip.Core) *Handlers {
	return &Handlers{
		trip: trip,
	}
}

// Create adds a new trip to the system.
func (h *Handlers) Create(ctx context.Context, c *gin.Context) error {
	var app AppNewTrip
	// Validate the request.
	if err := web.Decode(c, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	nc, err := toCoreNewTrip(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	fmt.Println("handlers: trip: create: nc:", nc)
	trip, err := h.trip.Create(ctx, nc)
	if err != nil {
		return fmt.Errorf("create: usr[%+v]: %w", trip, err)
	}

	return web.Respond(ctx, c.Writer, toAppTrip(trip), http.StatusCreated)
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

	trips, err := h.trip.QueryAll(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	items := make([]AppTrip, len(trips))
	for i, trip := range trips {
		items[i] = toAppTrip(trip)
	}

	total, err := h.trip.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, c.Writer, paging.NewResponse(items, total, page.Number, page.RowsPerPage), http.StatusOK)
}

// Query returns a trip that matches the specified ID in the session.
func (h *Handlers) QueryByUserID(ctx context.Context, c *gin.Context) error {
	id := auth.GetUserID(ctx)

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

	trip, err := h.trip.QueryByUserID(ctx, id, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: id[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, c.Writer, paging.NewResponse(trip, len(trip), page.Number, page.RowsPerPage), http.StatusOK)
}

// QueryByID returns a trip by its ID.
func (h *Handlers) QueryByID(ctx context.Context, c *gin.Context) error {
	id := c.Param("id")

	usr, err := h.trip.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: id[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, c.Writer, toAppTrip(usr), http.StatusOK)
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
