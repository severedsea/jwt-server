package redis

import (
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
)

// New returns a redis client
func New() (redis.Cmdable, error) {
	opt, err := redis.ParseURL(redisURLFromEnv())
	if err != nil {
		return nil, err
	}

	opt.MaxRetries = 3

	return redis.NewClient(opt), nil
}

// redisURLFromEnv constructs the redis URL from the REDIS_* env vars
func redisURLFromEnv() string {
	var sb strings.Builder

	// Scheme
	if scheme := os.Getenv("REDIS_SCHEME"); scheme != "" {
		sb.WriteString(scheme)
		sb.WriteString("://")
	}

	// Credentials
	if user, password := os.Getenv("REDIS_USER"), os.Getenv("REDIS_PWD"); user != "" || password != "" {
		sb.WriteString(user)
		sb.WriteString(":")
		sb.WriteString(password)
		sb.WriteString("@")
	}

	// Host
	sb.WriteString(os.Getenv("REDIS_HOST"))

	// Port
	if port := os.Getenv("REDIS_PORT"); port != "" {
		sb.WriteString(":")
		sb.WriteString(os.Getenv("REDIS_PORT"))
	}

	return sb.String()
}
