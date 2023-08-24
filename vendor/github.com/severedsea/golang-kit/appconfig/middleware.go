package appconfig

import (
	"net/http"

	"github.com/severedsea/golang-kit/web"
	"github.com/severedsea/golang-kit/web/middleware"
)

var ErrLoad = &web.Error{Status: http.StatusInternalServerError, Code: "appconfig_load"}

// Middleware serves as a middleware that loads the app config into the context
func Middleware() middleware.Adapter {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx, err := LoadFromEnv(ctx)
			if err != nil {
				web.RespondJSON(ctx, w, web.NewError(ErrLoad, err.Error()), nil)

				return
			}

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
