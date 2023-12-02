// Package authgrp maintains the group of handlers for user access.
package authgrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TSMC-Uber/server/business/core/auth"
	"github.com/TSMC-Uber/server/business/core/user"
	webauth "github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/business/web/v1/response"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/gin-gonic/gin"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	auth *auth.Core
	user *user.Core
}

// New constructs a handlers for route access.
func New(auth *auth.Core, user *user.Core) *Handlers {
	return &Handlers{
		auth: auth,
		user: user,
	}
}

// @Summary login
// @Schemes
// @Description Login will add a user if they do not exist or update them if they do.
// @Tags auth
// @Accept json
// @Produce json
// @Param id_token header string true "ID Token"
// @Success 201 {object} AppUser "User successfully logged in"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /auth/login [post]
// Login will add a user if they do not exist or update them if they do.
func (h *Handlers) Login(ctx context.Context, c *gin.Context) error {
	idToken := webauth.GetIDToken(ctx)

	tokenInfo, err := h.auth.ParseIDToken(idToken)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	nc, err := toCoreNewUser(tokenInfo)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	usr, err := h.user.UpsertByGoogleID(ctx, tokenInfo.Sub, nc)
	if err != nil {
		return fmt.Errorf("upsert: usr[%+v]: %w", usr, err)
	}

	// generate token
	sessionToken, err := h.auth.GenerateSessionToken(usr.ID, 3600)
	if err != nil {
		return fmt.Errorf("generate session token: %w", err)
	}

	// store session token to redis
	err = h.auth.SetSessionToken(ctx, *sessionToken, usr)
	if err != nil {
		return fmt.Errorf("set session token: %w", err)
	}
	// set cookie
	c.SetCookie("token", sessionToken.Plaintext, 3600, "/", "localhost", false, true)
	return web.Respond(ctx, c.Writer, toAppUser(usr), http.StatusCreated)
}

// Logout updates a user in the system.
func (h *Handlers) Logout(ctx context.Context, c *gin.Context) error {
	sessionToken := webauth.GetSessionToken(ctx)
	// remove session token from redis
	err := h.auth.RemoveSessionToken(ctx, sessionToken)
	if err != nil {
		return fmt.Errorf("remove session token: %w", err)
	}

	// remove cookie
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	return nil
}
