package mid

import (
	"context"

	"github.com/TSMC-Uber/server/business/web/metrics"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/gin-gonic/gin"
)

// Metrics updates program counters.
func Metrics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, c *gin.Context) error {
			ctx = metrics.Set(ctx)

			err := handler(ctx, c)

			metrics.AddRequests(ctx)
			metrics.AddGoroutines(ctx)

			if err != nil {
				metrics.AddErrors(ctx)
			}

			return err
		}

		return h
	}

	return m
}
