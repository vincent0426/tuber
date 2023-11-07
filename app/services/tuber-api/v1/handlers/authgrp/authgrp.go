// Package authgrp maintains the group of handlers for user access.
package authgrp

import (
	"context"

	"github.com/TSMC-Uber/server/business/core/auth"
	"github.com/gin-gonic/gin"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	auth *auth.Core
}

// New constructs a handlers for route access.
func New(auth *auth.Core) *Handlers {
	return &Handlers{
		auth: auth,
	}
}

// Login adds a new user to the system.
func (h *Handlers) Login(ctx context.Context, c *gin.Context) error {

	return nil
}

// Logout updates a user in the system.
func (h *Handlers) Logout(ctx context.Context, c *gin.Context) error {
	return nil
}
