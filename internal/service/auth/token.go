package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	jwtgo "github.com/golang-jwt/jwt/v4"
	"github.com/severedsea/jwt-server/internal/pkg/jwt"
)

const (
	tokenTypeBearer     TokenType = "Bearer"
	tokenExpiryDuration           = time.Duration(20) * time.Minute
)

// TokenType is the enum for exemption
type TokenType string

// IsValid checks is the value is in the enum list
func (e TokenType) IsValid() bool {
	return e == tokenTypeBearer
}

// String returns enum in string
func (e TokenType) String() string {
	return string(e)
}

// Token contains the token information
type Token struct {
	AccessToken string
	TokenType   TokenType
	ExpiresIn   int
	ExpiresAt   time.Time
	Scope       string
}

// Claims is the claims for the JWT
type Claims struct {
	jwtgo.RegisteredClaims
	// Add custom claims here if necessary
}

/*
redisValue is the value for storing token related data in redis

	It will implement encoding.BinaryMarshaler and encoding.BinaryUnMarshaler so that go-redis can unmarshal it automatically

https://github.com/go-redis/redis/issues/739#issuecomment-418185046
*/
type redisValue struct {
	AccessToken string `redis:"AccessToken"`
}

func (v redisValue) MarshalBinary() ([]byte, error) {
	return json.Marshal(v)
}

func (v *redisValue) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &v)
}

func (s Service) GenerateToken(ctx context.Context, subject string) (Token, error) {
	// Generate claims
	c := Claims{
		RegisteredClaims: jwt.NewRegisteredClaims(subject, tokenExpiryDuration),
	}

	// Sign claims
	tokenString, err := jwt.Sign(c)
	if err != nil {
		return Token{}, err
	}

	// Save token to redis
	key := redisKey(subject)
	if err := s.redis.SetArgs(ctx, key,
		redisValue{
			AccessToken: tokenString,
		},
		// [20210827] ExpiresAt is only available on redis >= 6.2, we're using AWS Elasticache 6.0.5
		redis.SetArgs{TTL: tokenExpiryDuration},
	).Err(); err != nil {
		return Token{}, err
	}

	return Token{
		AccessToken: tokenString,
		ExpiresIn:   int(tokenExpiryDuration.Seconds()),
		ExpiresAt:   time.Unix(c.ExpiresAt.Unix(), 0),
		TokenType:   tokenTypeBearer,
	}, nil
}

// ParseToken validates and parses the token string
func (s Service) ParseToken(_ context.Context, tokenString string) (Claims, error) {
	c := Claims{}
	if err := jwt.Parse(tokenString, &c); err != nil {
		return Claims{}, err
	}

	return c, nil
}

// VerifyToken verifies the token against the one stored in redis
func (s Service) VerifyToken(ctx context.Context, tokenString, subject string) error {
	// check if token in redis is equal
	var v redisValue
	if err := s.redis.Get(ctx, redisKey(subject)).Scan(&v); err != nil {
		return jwt.ErrInvalidToken
	}

	if v.AccessToken != tokenString {
		return jwt.ErrInvalidToken
	}

	return nil
}

func redisKey(subject string) string {
	return "auth_" + subject
}
