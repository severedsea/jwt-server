package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/severedsea/golang-kit/logr"
	"github.com/severedsea/golang-kit/web"
)

// RequestLogger serves as a middleware that logs the start and end of each request, along with some useful data as logger fields
func RequestLogger() Adapter {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ctx := r.Context()

			// Parse request information
			requestURIparts := append(strings.SplitN(r.RequestURI, "?", 2), "") // `append` so we'd always have an array of 2 strings at least
			rctx := web.NewRequestContext(r)
			ctx = web.SetRequestContext(ctx, rctx)

			// Instantiate verbose logger
			logger := logr.DefaultLogger().
				WithField("request", rctx.RequestID).
				WithField("route", r.Method+" "+requestURIparts[0]).
				WithField("query", requestURIparts[1]).
				WithField("ip", r.RemoteAddr).
				WithField("referer", r.Referer()).
				WithField("agent", r.UserAgent())

				// Set logger into context
			ctx = logr.SetLogger(ctx, logger)

			logger.
				Infof("START")

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

			logger.
				WithField("duration", time.Since(start)).
				Infof("END")
		}
		return http.HandlerFunc(fn)
	}
}
