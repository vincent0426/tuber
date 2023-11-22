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
	"github.com/google/uuid"
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
	userID := auth.GetUserID(ctx)
	var app AppNewTrip
	// Validate the request.
	if err := web.Decode(c, &app); err != nil {
		fmt.Println("app:create:decode:err:", err)
		return response.NewError(err, http.StatusBadRequest)
	}

	nc, err := toCoreNewTrip(app)
	if err != nil {
		fmt.Println("app:create:toCoreNewTrip:err:", err)
		return response.NewError(err, http.StatusBadRequest)
	}
	nc.DriverID = userID

	trip, err := h.trip.Create(ctx, nc)
	if err != nil {
		fmt.Println("app:create:trip:create:err:", err)
		return fmt.Errorf("create: trip[%+v]: %w", trip, err)
	}

	return web.Respond(ctx, c.Writer, toAppTrip(trip), http.StatusCreated)
}

// Update updates a trip in the system.
func (h *Handlers) Update(ctx context.Context, c *gin.Context) error {
	userID := auth.GetUserID(ctx)

	var app AppUpdateTrip
	if err := web.Decode(c, &app); err != nil {
		return err
	}

	tripID := uuid.Must(uuid.Parse(c.Param("id")))

	qtrip, err := h.trip.QueryByID(ctx, tripID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: tripID[%s]: %w", tripID, err)
		}
	}

	// check if the user is the driver of the trip
	if qtrip.DriverID != userID {
		fmt.Println("app:update:trip:userID:", userID, "qtrip.DriverID:", qtrip.DriverID)
		return response.NewError(errors.New("user is not the driver of the trip"), http.StatusForbidden)
	}

	ut, err := toCoreUpdateTrip(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	trip, err := h.trip.Update(ctx, qtrip, ut)
	if err != nil {
		return fmt.Errorf("update: tripID[%s] ut[%+v]: %w", tripID, ut, err)
	}

	return web.Respond(ctx, c.Writer, toAppTrip(trip), http.StatusOK)
}

// Query returns a list of users with paging.
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

	trips, err := h.trip.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	items := make([]AppTripView, len(trips))
	for i, trip := range trips {
		items[i] = toAppTripView(trip)
	}

	total, err := h.trip.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, c.Writer, paging.NewResponse(items, total, page.Number, page.RowsPerPage), http.StatusOK)
}

// QueryMyTrip returns a trip that the user is in.
func (h *Handlers) QueryMyTrip(ctx context.Context, c *gin.Context) error {
	id := auth.GetUserID(ctx)

	page, err := paging.ParseRequest(c.Request)
	if err != nil {
		return err
	}

	filter, err := parseFilterByUser(c.Request)
	if err != nil {
		return err
	}

	orderBy, err := parseOrder(c.Request)
	if err != nil {
		return err
	}

	qtrip, err := h.trip.QueryMyTrip(ctx, id, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		switch {
		case errors.Is(err, trip.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: id[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, c.Writer, paging.NewResponse(qtrip, len(qtrip), page.Number, page.RowsPerPage), http.StatusOK)
}

// QueryByID returns a trip by its ID.
func (h *Handlers) QueryByID(ctx context.Context, c *gin.Context) error {
	id := c.Param("id")

	qtrip, err := h.trip.QueryByID(ctx, uuid.Must(uuid.Parse(id)))
	if err != nil {
		switch {
		case errors.Is(err, trip.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: id[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, c.Writer, toAppTripView(qtrip), http.StatusOK)
}

func (h *Handlers) Join(ctx context.Context, c *gin.Context) error {
	userID := auth.GetUserID(ctx)
	var app AppNewTripPassenger
	// Validate the request.
	if err := web.Decode(c, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	ntp, err := toCoreNewTripPassenger(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	ntp.PassengerID = userID

	tripPassenger, err := h.trip.Join(ctx, ntp)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return response.NewError(err, http.StatusConflict)
		}
		return fmt.Errorf("join: tripPassenger[%+v]: %w", tripPassenger, err)
	}

	return web.Respond(ctx, c.Writer, toAppTripPassenger(tripPassenger), http.StatusCreated)
}
