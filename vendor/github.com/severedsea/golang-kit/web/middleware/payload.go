package middleware

import (
	"net/http"
)

// LimitPayload serves as a middleware that limits request payload size to the maximum configured size
func LimitPayload(maxbytes int64) Adapter {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxbytes)

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
