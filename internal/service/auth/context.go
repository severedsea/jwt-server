package auth

import (
	"context"
	"os"
)

type contextKey string

const (
	claimsContextKey = contextKey("jwt_claims")
)

// ClaimsFromContext returns Token details inside the context
func ClaimsFromContext(ctx context.Context) (Claims, error) {
	value, ok := ctx.Value(claimsContextKey).(*Claims)
	if value == nil || !ok {
		return Claims{}, ErrMissingContext
	}

	return *value, nil
}

func setClaimsContext(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, interface{}(claimsContextKey), &claims)
}

// =====================
// FOR USE IN TESTS ONLY
// =====================

// ProvidesClaims puts JWT claims into the context
// This is ONLY meant for tests.
func ProvidesClaims(ctx context.Context, claims Claims) context.Context {
	if os.Getenv("APP_ENV") != "test" {
		panic("NOT for use in APP_ENV=" + os.Getenv("APP_ENV"))
	}
	ctx = setClaimsContext(ctx, claims)

	return ctx
}
