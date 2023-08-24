package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

// Parse validates and parses the token string
func Parse(tokenString string, c jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, c, func(t *jwt.Token) (interface{}, error) {
		// Validate alg
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return verifyKey, nil
	})
	if err != nil {
		return ErrInvalidToken
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	mc, ok := c.(*jwt.RegisteredClaims)
	if !ok {
		return nil
	}

	// Best-effort validation of issuer
	if mc.Issuer != Issuer {
		return ErrInvalidToken
	}

	return nil
}
