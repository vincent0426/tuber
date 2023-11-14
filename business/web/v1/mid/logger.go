package mid

import (
	"context"
	"fmt"
	"time"

	"github.com/TSMC-Uber/server/foundation/logger"
	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/gin-gonic/gin"
)

func Logger(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, c *gin.Context) error {
			v := web.GetValues(ctx)
			r := c.Request
			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}

			log.Info(ctx, "request started", "method", r.Method, "path", path,
				"remoteaddr", r.RemoteAddr)

			err := handler(ctx, c)

			log.Info(ctx, "request completed", "method", r.Method, "path", path,
				"remoteaddr", r.RemoteAddr, "statuscode", v.StatusCode, "since", time.Since(v.Now))

			return err
		}

		return h
	}

	return m
}
