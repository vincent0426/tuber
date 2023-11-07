package mid

import (
	"context"

	aauth "github.com/TSMC-Uber/server/business/core/auth"
	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/gin-gonic/gin"
)

// auth idtoken
// auth google
func AuthGoogle(a *auth.Auth, authCore *aauth.Core) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, c *gin.Context) error {
			err := a.AuthGoogle(ctx, c.Request.Header.Get("id_token"))
			if err != nil {
				return auth.NewAuthError("authenticate google: failed: %s", err)
			}

			return handler(ctx, c)
		}

		return h
	}

	return m
}
