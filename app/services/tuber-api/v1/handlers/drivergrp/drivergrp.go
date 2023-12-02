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

// @Summary create a new driver
// @Schemes
// @Description Create will add a driver if they do not exist or update them if they do.
// @Tags driver
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param body body AppNewDriver true "New Driver"
// @Success 201 {object} AppDriver "Driver successfully created"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /drivers [post]
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

// @Summary get all drivers
// @Schemes
// @Description QueryAll will query drivers
// @Tags driver
// @Accept json
// @Produce json
// @Success 200 {object} AppDriver "Drivers successfully queried"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /drivers [get]
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

	drivers, err := h.driver.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	items := make([]AppDriver, len(drivers))
	for i, driver := range drivers {
		items[i] = toAppDriver(driver)
	}

	total, err := h.driver.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, c.Writer, paging.NewResponse(items, total, page.Number, page.RowsPerPage), http.StatusOK)
}

// @Summary get drivers by id
// @Schemes
// @Description QueryByID will query drivers by id
// @Tags driver
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} AppDriver "Drivers successfully queried"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /drivers/{id} [get]
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

// @Summary add a driver to favorite
// @Schemes
// @Description AddFavorite will add a driver to favorite
// @Tags driver
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param id path string true "ID"
// @Success 200 {object} AppDriver "Drivers successfully queried"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /favorite-drivers/{id} [post]
func (h *Handlers) AddFavorite(ctx context.Context, c *gin.Context) error {
	userID := auth.GetUserID(ctx)
	driverID := c.Param("id")

	if err := h.driver.AddFavorite(ctx, userID, driverID); err != nil {
		return fmt.Errorf("addfavorite: driverID[%s]: %w", driverID, err)
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("driver[%s] added to favorite", driverID),
	}

	return web.Respond(ctx, c.Writer, response, http.StatusOK)
}

// query favorite drivers
// @Summary get favorite drivers
// @Schemes
// @Description QueryFavorite will query favorite drivers
// @Tags driver
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Success 200 {object} AppDriver "Drivers successfully queried"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /favorite-drivers [get]
func (h *Handlers) QueryFavorite(ctx context.Context, c *gin.Context) error {
	userID := auth.GetUserID(ctx)

	page, err := paging.ParseRequest(c.Request)
	if err != nil {
		return err
	}

	filter, err := parseFilter(c.Request)
	if err != nil {
		return err
	}

	orderBy, err := parsefavoriteDriverOrder(c.Request)
	if err != nil {
		return err
	}

	favoriteDrivers, err := h.driver.QueryFavorite(ctx, userID, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	items := make([]AppFavoriteDriver, len(favoriteDrivers))
	for i, driver := range favoriteDrivers {
		items[i] = toAppFavoriteDriver(driver)
	}

	total, err := h.driver.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, c.Writer, paging.NewResponse(items, total, page.Number, page.RowsPerPage), http.StatusOK)
}
