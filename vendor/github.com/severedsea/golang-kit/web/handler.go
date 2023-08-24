package web

import (
	"net/http"
)

// WrapHandler wraps the web.HandlerFunc to standard http.HandlerFunc with error handling
func WrapHandler(h HandlerFunc) http.HandlerFunc {
	wh := Handler{H: h}

	return wh.ServeHTTP
}

// HandlerFunc is a http.HandlerFunc variant that returns error
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// Handler is a http.Handler implementation that handles HandlerFunc
type Handler struct {
	H HandlerFunc
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.H(w, r); err != nil {
		ctx := r.Context()
		RespondJSON(ctx, w, err, nil)
	}
}
