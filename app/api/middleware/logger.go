package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/danipurwadi/internal-transfer-system/foundation/logger"
	"github.com/danipurwadi/internal-transfer-system/foundation/web"
)

func Logger(log *logger.Logger) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v := web.GetValues(ctx)

			log.Info(ctx, "request started", "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Info(ctx, "request completed", "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr,
				"statuscode", v.StatusCode, "sinceInMs", time.Since(v.Now).Milliseconds())
			return err
		}
		return h
	}
	return m
}
