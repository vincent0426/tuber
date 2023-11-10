// Package authgrp maintains the group of handlers for user access.
package authgrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TSMC-Uber/server/business/core/auth"
	"github.com/TSMC-Uber/server/business/core/user"
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

// Login adds a new user to the system.
func (h *Handlers) Login(ctx context.Context, c *gin.Context) error {
	idToken := c.Request.Header.Get("id_token")

	tokenInfo, err := h.auth.GetTokenInfo(idToken)
	if err != nil {
		return fmt.Errorf("get token info: %w", err)
	}

	nc, err := toCoreNewUser(tokenInfo)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	fmt.Println(nc)

	// usr, err := h.user.UpsertByGoogleID(ctx, tokenInfo.Id, nc)
	// if err != nil {
	// 	return fmt.Errorf("upsert user: %w", err)
	// }

	// // generate token
	// sessionToken, err := webauth.GenerateSessionToken(usr.ID, 3600)
	// if err != nil {
	// 	return fmt.Errorf("generate token: %w", err)
	// }

	// // insert token to db
	// coreSessionToken := toCoreSessionToken(sessionToken)
	// err = h.auth.UpsertSessionToken(ctx, coreSessionToken)
	// if err != nil {
	// 	return fmt.Errorf("upsert token: %w", err)
	// }

	// set cookie
	// c.SetCookie("token", sessionToken.Plaintext, 3600, "/", "localhost", false, true)
	return web.Respond(ctx, c.Writer, "login success", http.StatusOK)
}

// Logout updates a user in the system.
func (h *Handlers) Logout(ctx context.Context, c *gin.Context) error {
	return nil
}
