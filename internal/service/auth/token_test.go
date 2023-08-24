package auth

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/severedsea/golang-kit/appconfig"
	"github.com/severedsea/jwt-server/internal/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestToken_IsValid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		given    TokenType
		expected bool
	}{
		{TokenType("invalid"), false},
		{tokenTypeBearer, true},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.given.String(), func(t *testing.T) {
			t.Parallel()
			// Given:

			// When:
			actual := tc.given.IsValid()

			// Then:
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestParseToken(t *testing.T) {
	t.Parallel()

	// Given:
	ctx := context.Background()
	ctx, err := appconfig.LoadFromEnv(ctx)
	assert.NoError(t, err)
	subject := "123"

	// Mocks:
	mockRds := &mockRedis{}
	mockRds.On("SetArgs", mock.Anything, redisKey(subject),
		mock.AnythingOfType("redisValue"), mock.AnythingOfType("redis.SetArgs")).
		Return(redis.NewStatusResult("", nil))

	s := New(mockRds)
	// gen a new Token
	exp, err := s.GenerateToken(ctx, subject)
	assert.NoError(t, err)

	// When:
	act, err := s.ParseToken(ctx, exp.AccessToken)

	// Then:
	assert.NoError(t, err)

	assert.Equal(t, "jwt-server", act.Issuer)
	assert.Equal(t, subject, act.Subject)
	assert.Equal(t, exp.ExpiresAt.UnixNano(), act.ExpiresAt.UnixNano())

	iat := exp.ExpiresAt.Add(-20 * time.Minute)
	assert.Equal(t, iat.UnixNano(), act.IssuedAt.UnixNano())
}

func TestParseToken_Error(t *testing.T) {
	t.Parallel()

	// Given:
	ctx := context.Background()

	// When:
	s := New(nil)
	act, err := s.ParseToken(ctx, "INVALID_ACCESS_TOKEN")

	// Then:
	assert.Error(t, err)
	assert.Empty(t, act)
	assert.ErrorIs(t, err, jwt.ErrInvalidToken)
}

func TestVerifyToken(t *testing.T) {
	t.Parallel()

	// Given:
	ctx := context.Background()
	subject := "4321"
	expClaims := Claims{
		RegisteredClaims: jwt.NewRegisteredClaims(subject, time.Hour),
	}
	tokenString, err := jwt.Sign(expClaims)
	assert.NoError(t, err)

	b, err := json.Marshal(redisValue{AccessToken: tokenString})
	assert.NoError(t, err)

	// Mocks:
	mockRds := &mockRedis{}
	mockRds.On("Get", mock.Anything, redisKey(subject)).
		Return(redis.NewStringResult(string(b), nil))

	// When:
	s := New(mockRds)
	err = s.VerifyToken(ctx, tokenString, subject)

	// Then:
	assert.NoError(t, err)

	// Assert mocks call
	mockRds.AssertNumberOfCalls(t, "Get", 1)
}

func TestVerifyToken_Error(t *testing.T) {
	t.Parallel()

	subject := "1234"
	expClaims := Claims{
		RegisteredClaims: jwt.NewRegisteredClaims(subject, time.Hour),
	}
	tokenString, err := jwt.Sign(expClaims)
	assert.NoError(t, err)

	testCases := []struct {
		desc     string
		given    string
		expCalls int
		exp      error
		mock     func(rds *mockRedis)
	}{
		{
			desc:     "Redis error",
			given:    tokenString,
			expCalls: 1,
			exp:      jwt.ErrInvalidToken,
			mock: func(rds *mockRedis) {
				rds.On("Get", mock.Anything, redisKey(subject)).
					Return(redis.NewStringResult("", redis.ErrClosed))
			},
		},
		{
			desc:     "JWT Token is not equal",
			given:    tokenString,
			expCalls: 1,
			exp:      jwt.ErrInvalidToken,
			mock: func(rds *mockRedis) {
				rds.On("Get", mock.Anything, redisKey(subject)).
					Return(redis.NewStringResult("NOT THE SAME TOKEN", nil))
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
			tc.mock(mockRds)

			// When:
			s := New(mockRds)
			err := s.VerifyToken(ctx, tc.given, subject)

			// Then:
			assert.Error(t, err)
			assert.ErrorIs(t, err, tc.exp)
			assert.Equal(t, tc.exp, err)

			// Assert mocks call
			mockRds.AssertNumberOfCalls(t, "Get", tc.expCalls)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	// Given:
	ctx := context.Background()
	ctx, err := appconfig.LoadFromEnv(ctx)
	assert.NoError(t, err)

	subject := "123"

	mockRds := &mockRedis{}
	s := New(mockRds)

	expClaims := Claims{
		RegisteredClaims: jwt.NewRegisteredClaims(subject, tokenExpiryDuration),
	}
	expTokenString, err := jwt.Sign(expClaims)
	assert.NoError(t, err)

	// Mocks:
	mockRds.On("SetArgs", mock.Anything, redisKey(subject),
		redisValue{
			AccessToken: expTokenString,
		},
		redis.SetArgs{TTL: tokenExpiryDuration}).
		Return(redis.NewStatusResult("", nil))

	// When:
	act, err := s.GenerateToken(ctx, subject)

	// Then:
	assert.NoError(t, err)
	assert.Equal(t, Token{
		AccessToken: expTokenString,
		TokenType:   tokenTypeBearer,
		ExpiresIn:   int(tokenExpiryDuration.Seconds()),
		ExpiresAt:   time.Unix(expClaims.ExpiresAt.Unix(), 0),
	}, act)

	// Assert mocks call
	mockRds.AssertNumberOfCalls(t, "SetArgs", 1)
}

func TestGenerateToken_Error(t *testing.T) {
	mockRds := &mockRedis{}
	subject := "ID_NO"

	testCases := []struct {
		desc  string
		mocks func(ctx context.Context) context.Context
		exp   error
	}{
		{
			desc: "redis set error",
			mocks: func(ctx context.Context) context.Context {
				var err error
				ctx, err = appconfig.LoadFromEnv(ctx)
				assert.NoError(t, err)

				mockRds.On("SetArgs", mock.Anything, redisKey(subject),
					mock.AnythingOfType("redisValue"), mock.AnythingOfType("redis.SetArgs")).
					Return(redis.NewStatusResult("", redis.Nil))

				return ctx
			},
			exp: redis.Nil,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			ctx := context.Background()

			// Mocks:
			ctx = tc.mocks(ctx)

			// When:
			s := New(mockRds)
			act, err := s.GenerateToken(ctx, subject)

			// Then:
			assert.Error(t, err)
			assert.Empty(t, act)
			assert.ErrorIs(t, err, tc.exp)
		})
	}
}
