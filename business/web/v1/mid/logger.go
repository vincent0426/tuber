package mid

import (
	"context"
	"fmt"
	"time"

	"github.com/TSMC-Uber/server/foundation/web"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, c *gin.Context) error {
			v := web.GetValues(ctx)
			r := c.Request
			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}

			log.Infow("request started", "trace_id", v.TraceID, "method", r.Method, "path", path,
				"remoteaddr", r.RemoteAddr, "useragent", r.UserAgent())

			err := handler(ctx, c)

			log.Infow("request completed", "trace_id", v.TraceID, "method", r.Method, "path", path,
				"remoteaddr", r.RemoteAddr, "statuscode", v.StatusCode, "since", time.Since(v.Now))

			return err
		}

		return h
	}

	return m
}
