package middleware

import (
	"fmt"
	"net/http"
)

// MaxAge serves as a middleware that sets `max-age` value to `Cache-Control` HTTP header
func MaxAge(seconds int) Adapter {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Set max-age HTTP header value
			w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", seconds))

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
