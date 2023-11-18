// Package authgrp maintains the group of handlers for user access.
package authgrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TSMC-Uber/server/business/core/auth"
	"github.com/TSMC-Uber/server/business/core/user"
	webauth "github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/business/web/v1/mid"
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

// Login will add a user if they do not exist or update them if they do.
func (h *Handlers) Login(ctx context.Context, c *gin.Context) error {
	idToken := webauth.GetIDToken(ctx)

	tokenInfo, err := h.auth.ParseIDToken(idToken)
	if err != nil {
		return mid.WrapError(fmt.Errorf("parse id token: %w", err))
	}

	nc, err := toCoreNewUser(tokenInfo)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	usr, err := h.user.UpsertByGoogleID(ctx, tokenInfo.Sub, nc)
	if err != nil {
		return mid.WrapError(fmt.Errorf("upsert user: %w", err))
	}

	// generate token
	sessionToken, err := h.auth.GenerateSessionToken(usr.ID, 3600)
	if err != nil {
		return mid.WrapError(fmt.Errorf("generate session token: %w", err))
	}

	// store session token to redis
	err = h.auth.SetSessionToken(ctx, *sessionToken, usr)
	if err != nil {
		return mid.WrapError(fmt.Errorf("set session token: %w", err))
	}
	// set cookie
	c.SetCookie("token", sessionToken.Plaintext, 3600, "/", "localhost", false, true)
	return web.Respond(ctx, c.Writer, fmt.Sprintf("Welcome %s", usr.Name), http.StatusOK)
}

// Logout updates a user in the system.
func (h *Handlers) Logout(ctx context.Context, c *gin.Context) error {
	return nil
}
