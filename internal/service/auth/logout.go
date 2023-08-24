package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/severedsea/golang-kit/logr"
	"github.com/severedsea/golang-kit/timex"
	"github.com/severedsea/golang-kit/web"
)

// Logout invalidates the access_token for the provided subject
func (s Service) Logout(ctx context.Context, subject string) error {
	logger := logr.GetLogger(ctx)
	startTime := timex.NowSGT()

	// retrieve value from redis
	var v redisValue
	key := redisKey(subject)
	if err := s.redis.Get(ctx, key).Scan(&v); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}

		return web.NewError(ErrRedis, err.Error())
	}

	// delete key in redis
	d, err := s.redis.Del(ctx, key).Result()
	if err != nil {
		return web.NewError(ErrRedis, err.Error())
	}
	if d < 1 {
		return web.NewError(ErrRedis, fmt.Sprintf("key %s was not deleted", key))
	}

	logger.
		WithField("duration", time.Since(startTime).Milliseconds()).
		Infof("logout successful")

	return nil
}
