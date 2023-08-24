package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
)

var realm = getRealm()

// BasicAuth serves as a middleware to handle basic authentication
// Reference: http://stackoverflow.com/a/39591234/136558
func BasicAuth(username string, password string) Adapter {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Basic auth disabled
			if username == "" && password == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Invalid auth
			u, p, ok := r.BasicAuth()
			if !ok ||
				subtle.ConstantTimeCompare([]byte(u), []byte(username)) != 1 ||
				subtle.ConstantTimeCompare([]byte(p), []byte(password)) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorised.\n"))
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func getRealm() string {
	if os.Getenv("HTTP_BASIC_REALM") != "" {
		return os.Getenv("HTTP_BASIC_REALM")
	}
	if os.Getenv("APPNAME") != "" {
		return os.Getenv("APPNAME")
	}
	return ""
}
