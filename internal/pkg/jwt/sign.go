package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/severedsea/golang-kit/web"
)

// Sign signs the JWT claims and returns the JWT string
func Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign claims
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", web.NewError(ErrJWT, err.Error())
	}

	return tokenString, nil
}
