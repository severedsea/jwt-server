package auth

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// New creates a new Service struct
func New(rds redis.Cmdable) Service {
	return Service{
		redis: rds,
	}
}

// Service holds the methods for this package
type Service struct {
	redis redis.Cmdable
}

// TokenParser is the interface for the token parser
type TokenParser interface {
	ParseToken(ctx context.Context, tokenString string) (Claims, error)
}

type TokenVerifier interface {
	VerifyToken(ctx context.Context, tokenString, subject string) error
}

type TokenParserVerifier interface {
	TokenParser
	TokenVerifier
}
