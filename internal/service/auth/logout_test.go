package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogout(t *testing.T) {
	subject := "sub"
	key := redisKey(subject)

	testCases := []struct {
		desc    string
		mocks   func(r *mockRedis)
		asserts func(t *testing.T, r *mockRedis)
	}{
		{
			desc: "success",
			mocks: func(r *mockRedis) {
				r.On("Get", mock.Anything, key).
					Return(redis.NewStringResult(`{"AccessToken":"ACCESS_TOKEN"}`, nil)).
					On("Del", mock.Anything, []string{key}).
					Return(redis.NewIntResult(1, nil))
			},
			asserts: func(t *testing.T, r *mockRedis) {
				r.AssertNumberOfCalls(t, "Get", 1)
				r.AssertNumberOfCalls(t, "Del", 1)
			},
		},
		{
			desc: "redis.Nil error",
			mocks: func(r *mockRedis) {
				r.On("Get", mock.Anything, key).
					Return(redis.NewStringResult("", redis.Nil))
			},
			asserts: func(t *testing.T, r *mockRedis) {
				r.AssertNumberOfCalls(t, "Get", 1)
				r.AssertNotCalled(t, "Del")
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			ctx := context.Background()

			// Mocks:
			mockRds := &mockRedis{}
			tc.mocks(mockRds)

			// When:
			s := New(mockRds)
			err := s.Logout(ctx, subject)

			// Then:
			assert.NoError(t, err)

			// Assert mocks call
			tc.asserts(t, mockRds)
		})
	}
}

func TestLogout_Error(t *testing.T) {
	t.Parallel()

	subject := "sub"
	key := redisKey(subject)
	givenErr := errors.New("something happened")

	testCases := []struct {
		desc    string
		mocks   func(r *mockRedis)
		asserts func(t *testing.T, r *mockRedis)
		exp     error
	}{
		{
			desc: "redis Get error",
			mocks: func(r *mockRedis) {
				r.On("Get", mock.Anything, key).
					Return(redis.NewStringResult(``, givenErr))
			},
			asserts: func(t *testing.T, r *mockRedis) {
				r.AssertNumberOfCalls(t, "Get", 1)
				r.AssertNotCalled(t, "Del")
			},
			exp: ErrRedis,
		},
		{
			desc: "redis Del error",
			mocks: func(r *mockRedis) {
				r.On("Get", mock.Anything, key).
					Return(redis.NewStringResult(`{"AccessToken":"ACCESS_TOKEN"}`, nil)).
					On("Del", mock.Anything, []string{key}).
					Return(redis.NewIntResult(0, givenErr))
			},
			asserts: func(t *testing.T, r *mockRedis) {
				r.AssertNumberOfCalls(t, "Get", 1)
				r.AssertNumberOfCalls(t, "Del", 1)
			},
			exp: ErrRedis,
		},
		{
			desc: "redis Del error: key not deleted error",
			mocks: func(r *mockRedis) {
				r.On("Get", mock.Anything, key).
					Return(redis.NewStringResult(`{"AccessToken":"ACCESS_TOKEN"}`, nil)).
					On("Del", mock.Anything, []string{key}).
					Return(redis.NewIntResult(0, nil))
			},
			asserts: func(t *testing.T, r *mockRedis) {
				r.AssertNumberOfCalls(t, "Get", 1)
				r.AssertNumberOfCalls(t, "Del", 1)
			},
			exp: ErrRedis,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			// Given:
			ctx := context.Background()

			// Mocks:
			mockRds := &mockRedis{}
			tc.mocks(mockRds)

			// 	When:
			s := New(mockRds)
			err := s.Logout(ctx, subject)

			// Then:
			assert.Error(t, err)
			assert.ErrorIs(t, err, tc.exp)

			// Assert mocks call
			tc.asserts(t, mockRds)
		})
	}
}
