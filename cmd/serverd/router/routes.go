// Package router contains routing configuration for webizapid
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/severedsea/jwt-server/cmd/serverd/router/api"
)

// Handler returns the http handler that handles all requests
func Handler() http.Handler {
	r := chi.NewRouter()

	// Top-level middlewares
	r.Use(chimiddleware.Recoverer)

	// API routes
	r.Group(api.Router)

	return r
}
