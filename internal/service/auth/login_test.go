package auth

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	// Given:
	ctx := context.Background()
	subject := "sub"

	// Mocks:
	mockRds := &mockRedis{}
	mockRds.On("SetArgs", mock.Anything, redisKey(subject),
		mock.AnythingOfType("redisValue"), mock.AnythingOfType("redis.SetArgs")).
		Return(redis.NewStatusResult("", nil))

	// When:
	s := New(mockRds)
	act, err := s.Login(ctx, subject)

	// Then:
	assert.NoError(t, err)

	assert.Equal(t, tokenTypeBearer, act.TokenType)
	assert.NotEmpty(t, act.AccessToken)
	assert.NotEmpty(t, act.ExpiresIn)
	assert.NotEmpty(t, act.ExpiresAt)

	// Assert mocks call
	mockRds.AssertNumberOfCalls(t, "SetArgs", 1)
}
