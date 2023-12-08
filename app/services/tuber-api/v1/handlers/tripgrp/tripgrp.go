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
	user *user.Core
}

// New constructs a handlers for route access.
func New(trip *trip.Core, user *user.Core) *Handlers {
	return &Handlers{
		trip: trip,
		user: user,
	}
}

// @Summary create a new trip
// @Schemes
// @Description Create will add a trip
// @Tags trip
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param body body AppNewTrip true "New Trip"
// @Success 201 {object} AppTrip "Trip successfully created"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /trips [post]
func (h *Handlers) Create(ctx context.Context, c *gin.Context) error {
	userID := auth.GetUserID(ctx)
	var app AppNewTrip
	// Validate the request.
	if err := web.Decode(c, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	nc, err := toCoreNewTrip(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	nc.DriverID = userID

	trip, err := h.trip.Create(ctx, nc)
	if err != nil {
		return fmt.Errorf("create: trip[%+v]: %w", trip, err)
	}

	// get driver's email
	qdriver, err := h.user.QueryByID(ctx, trip.DriverID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: driverID[%s]: %w", trip.DriverID, err)
		}
	}

	err = h.trip.PubEventToMQ(ctx, nc.StartTime, qdriver.Email.Address)
	if err != nil {
		return fmt.Errorf("pub event to mq: %w", err)
	}

	return web.Respond(ctx, c.Writer, toAppTrip(trip), http.StatusCreated)
}

// @Summary update a trip
// @Schemes
// @Description Update will update a trip
// @Tags trip
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param id path string true "Trip ID"
// @Param body body AppUpdateTrip true "Update Trip"
// @Success 200 {object} AppTrip "Trip successfully updated"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /trips/{id} [put]
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

// @Summary get all trips
// @Schemes
// @Description Query will query trips
// @Tags trip
// @Accept json
// @Produce json
// @Success 200 {object} AppTrip "query trips"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /trips [get]
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

// @Summary get my trips
// @Schemes
// @Description QueryMyTrip will query trips by user
// @Tags trip
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param status query string false "Status"
// @Param is_driver query bool false "Is Driver"
// @Success 200 {object} AppTrip "query trips"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /trips/my [get]
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

// @Summary get trips by id
// @Schemes
// @Description QueryByID will query trips by id
// @Tags trip
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} AppTrip "query trips"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /trips/{id} [get]
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

// @Summary join a trip
// @Schemes
// @Description Join will join a trip
// @Tags trip
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param id path string true "Trip ID"
// @Success 200 {object} AppTrip "Trip successfully joined"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /trips/{id}/join [post]
func (h *Handlers) Join(ctx context.Context, c *gin.Context) error {
	tripID := uuid.Must(uuid.Parse(c.Param("id")))
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

	tripPassenger, err := h.trip.Join(ctx, tripID, ntp)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return response.NewError(err, http.StatusConflict)
		}
		return fmt.Errorf("join: tripPassenger[%+v]: %w", tripPassenger, err)
	}

	return web.Respond(ctx, c.Writer, toAppTripPassenger(tripPassenger), http.StatusCreated)
}

// @Summary get all passengers of a trip
// @Schemes
// @Description QueryPassengers will query passengers of a trip
// @Tags trip
// @Accept json
// @Produce json
// @Param id path string true "Trip ID"
// @Param token header string true "Token"
// @Success 200 {object} AppTripDetails "query passengers of a trip"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /trips/{id}/passengers [get]
func (h *Handlers) QueryPassengers(ctx context.Context, c *gin.Context) error {
	tripID := c.Param("id")

	tripDetails, err := h.trip.QueryPassengers(ctx, uuid.Must(uuid.Parse(tripID)))
	if err != nil {
		switch {
		case errors.Is(err, trip.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querypassengers: tripID[%s]: %w", tripID, err)
		}
	}

	return web.Respond(ctx, c.Writer, toAppTripDetails(tripDetails), http.StatusOK)
}

// @Summary create a new rating for a trip
// @Schemes
// @Description Create will add a rating
// @Tags trip
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param body body AppNewRating true "New Rating"
// @Success 201 {object} AppRating "Rating successfully created"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /trips/{id}/rating [post]
func (h *Handlers) CreateRating(ctx context.Context, c *gin.Context) error {
	tripID := uuid.Must(uuid.Parse(c.Param("id")))
	userID := auth.GetUserID(ctx)
	var app AppNewRating
	// Validate the request.
	if err := web.Decode(c, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	nr, err := toCoreNewRating(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	nr.CommenterID = userID

	rating, err := h.trip.CreateRating(ctx, tripID, nr)
	if err != nil {
		return fmt.Errorf("create: rating[%+v]: %w", rating, err)
	}

	return web.Respond(ctx, c.Writer, toAppRating(rating), http.StatusCreated)
}
