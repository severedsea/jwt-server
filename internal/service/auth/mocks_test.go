package auth

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
)

// mockRedis is the mock redis
type mockRedis struct {
	mock.Mock
	redis.Cmdable
}

func (m *mockRedis) SetArgs(ctx context.Context, key string, value interface{}, a redis.SetArgs) *redis.StatusCmd {
	args := m.Called(ctx, key, value, a)

	return args.Get(0).(*redis.StatusCmd)
}

func (m *mockRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)

	return args.Get(0).(*redis.StringCmd)
}

func (m *mockRedis) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)

	return args.Get(0).(*redis.IntCmd)
}
