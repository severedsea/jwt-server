package appconfig

import (
	"context"
	"errors"
)

type contextKey string

var (
	ErrMissingConfig    = errors.New("missing app config in context")
	appconfigContextKey = contextKey("appconfig")
)

type Config struct {
	JWTExpiryInMinutes int `json:"JWT_EXPIRY_MIN"`
}

// Get returns the app config from the context provided
func Get(ctx context.Context) (Config, error) {
	cfg, ok := ctx.Value(appconfigContextKey).(Config)
	if !ok {
		return Config{}, ErrMissingConfig
	}

	return cfg, nil
}

// Set sets the provided config into the context
func Set(ctx context.Context, value Config) context.Context {
	return context.WithValue(ctx, appconfigContextKey, value)
}
