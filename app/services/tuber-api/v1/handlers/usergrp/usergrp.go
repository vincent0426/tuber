// Package usergrp maintains the group of handlers for user access.
package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/TSMC-Uber/server/business/web/v1/paging"
	"github.com/TSMC-Uber/server/business/web/v1/response"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	user *user.Core
}

// New constructs a handlers for route access.
func New(user *user.Core) *Handlers {
	return &Handlers{
		user: user,
	}
}

// @Summary create a new user
// @Schemes
// @Description Create will add a user if they do not exist or update them if they do.
// @Tags user
// @Accept json
// @Produce json
// @Param body body AppNewUser true "New User"
// @Success 201 {object} AppUser "User successfully created"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /users [post]
func (h *Handlers) Create(ctx context.Context, c *gin.Context) error {
	var app AppNewUser
	// Validate the request.
	if err := web.Decode(c, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	nc, err := toCoreNewUser(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return response.NewError(err, http.StatusConflict)
		}
		return fmt.Errorf("create: usr[%+v]: %w", usr, err)
	}

	return web.Respond(ctx, c.Writer, toAppUser(usr), http.StatusCreated)
}

// @Summary update a user
// @Schemes
// @Description Update will update a user
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param body body AppUpdateUser true "Update User"
// @Param token header string true "Token"
// @Success 200 {object} AppUser "User successfully updated"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /users/{id} [put]
func (h *Handlers) Update(ctx context.Context, c *gin.Context) error {
	var app AppUpdateUser
	if err := web.Decode(c, &app); err != nil {
		return err
	}

	// userID := auth.GetUserID(ctx)
	// temp get from url "/users/:id" -> uuid
	userID := uuid.Must(uuid.Parse(c.Param("id")))

	usr, err := h.user.QueryByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: userID[%s]: %w", userID, err)
		}
	}

	uu, err := toCoreUpdateUser(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	usr, err = h.user.Update(ctx, usr, uu)
	if err != nil {
		return fmt.Errorf("update: userID[%s] uu[%+v]: %w", userID, uu, err)
	}

	return web.Respond(ctx, c.Writer, toAppUser(usr), http.StatusOK)
}

// @Summary delete a user
// @Schemes
// @Description Delete will delete a user
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param token header string true "Token"
// @Success 204 "No Content"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /users/{id} [delete]
func (h *Handlers) Delete(ctx context.Context, c *gin.Context) error {
	// userID := auth.GetUserID(ctx)
	userID := uuid.Must(uuid.Parse(c.Param("id")))

	usr, err := h.user.QueryByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return web.Respond(ctx, c.Writer, nil, http.StatusNoContent)
		default:
			return fmt.Errorf("querybyid: userID[%s]: %w", userID, err)
		}
	}

	if err := h.user.Delete(ctx, usr); err != nil {
		return fmt.Errorf("delete: userID[%s]: %w", userID, err)
	}

	return web.Respond(ctx, c.Writer, nil, http.StatusNoContent)
}

// @Summary get all users
// @Schemes
// @Description Query will query users
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} AppUser "query users"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /users [get]
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

	users, err := h.user.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	items := make([]AppUser, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	total, err := h.user.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, c.Writer, paging.NewResponse(items, total, page.Number, page.RowsPerPage), http.StatusOK)
}

// @Summary get users by id
// @Schemes
// @Description QueryByID will query users by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} AppUser "query users"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /users/{id} [get]
func (h *Handlers) QueryByID(ctx context.Context, c *gin.Context) error {
	id := uuid.Must(uuid.Parse(c.Param("id")))

	usr, err := h.user.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: id[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, c.Writer, toAppUser(usr), http.StatusOK)
}
