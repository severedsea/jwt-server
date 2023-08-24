// Package api contains all api handlers
package api

import (
	"github.com/go-chi/chi/v5"
	v1 "github.com/severedsea/jwt-server/cmd/serverd/router/api/v1"
)

// Router registers handlers to the router provided in the argument
func Router(r chi.Router) {
	// Middlewares

	// Versioned routes
	r.Group(v1.Router)
}
