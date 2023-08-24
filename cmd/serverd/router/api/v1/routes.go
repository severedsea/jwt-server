// Package v1 contains webizapid v1 API handlers
package v1

import (
	"log"

	"github.com/go-chi/chi/v5"
	goredis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/severedsea/jwt-server/internal/pkg/redis"
	"github.com/severedsea/jwt-server/internal/service/auth"
)

var (
	redisClient goredis.Cmdable
)

func init() {
	var err error
	redisClient, err = redis.New()
	if err != nil {
		log.Fatalf("%s", errors.Wrap(err, "redis"))
	}
}

// Router registers handlers to the router provided in the argument
func Router(r chi.Router) {
	r.Group(public)
	r.Group(authenticated)
}

func public(r chi.Router) {

	authSvc := auth.New(redisClient)
	a := NewAuthHandler(authSvc)

	r.Get("/v1/login", a.Login())

}

func authenticated(r chi.Router) {
	authSvc := auth.New(redisClient)

	// Middlewares
	// Authentication middleware - Parses the header and validates the token
	r.Use(auth.Middleware(authSvc))

	a := NewAuthHandler(authSvc)
	r.Post("/v1/verify", a.Verify())
	r.Get("/v1/logout", a.Logout())
}
