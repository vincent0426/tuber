package mid

import (
	"context"

	aauth "github.com/TSMC-Uber/server/business/core/auth"
	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/gin-gonic/gin"
)

// Authenticate validates a JWT from the `Authorization` header.
func Authenticate(a *auth.Auth, authCore *aauth.Core) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, c *gin.Context) error {
			userID, err := a.Authenticate(ctx, c.Request.Header.Get("Authorization"), authCore)
			if err != nil {
				//
				return auth.NewAuthError("authenticate: failed: %s", err)
			}
			ctx = auth.SetUserID(ctx, userID)

			return handler(ctx, c)
		}

		return h
	}

	return m
}

func AuthGoogle(a *auth.Auth, authCore *aauth.Core) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, c *gin.Context) error {
			idToken := c.Request.Header.Get("id_token")
			err := a.ValidateIDToken(idToken)
			if err != nil {
				return auth.NewAuthError("google authenticate: failed: %s", err)
			}
			ctx = auth.SetIDToken(ctx, idToken)

			return handler(ctx, c)
		}

		return h
	}

	return m
}

// Authorize validates that an authenticated user has at least one role from a
// specified list. This method constructs the actual function that is used.
// func Authorize(a *auth.Auth, rule string) web.Middleware {
// 	m := func(handler web.Handler) web.Handler {
// 		h := func(ctx context.Context, c *gin.Context) error {
// 			claims := auth.GetClaims(ctx)
// 			if claims.Subject == "" {
// 				return auth.NewAuthError("authorize: you are not authorized for that action, no claims")
// 			}

// 			if err := a.Authorize(ctx, claims, rule); err != nil {
// 				return auth.NewAuthError("authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", nil, rule, err)
// 			}

// 			return handler(ctx, c)
// 		}

// 		return h
// 	}

// 	return m
// }
