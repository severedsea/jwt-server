package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/severedsea/golang-kit/timex"
)

const (
	// Issuer is the constant for the jwt-server issuer value
	Issuer = "jwt-server"
)

// NewRegisteredClaims returns a new standard claims with the basic claims populated
func NewRegisteredClaims(subject string, tokenExpiryDuration time.Duration) jwt.RegisteredClaims {
	now := timex.NowSGT()
	expiresAt := now.Add(tokenExpiryDuration)

	return jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Issuer:    Issuer,
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(now),
	}
}
