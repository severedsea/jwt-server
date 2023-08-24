package jwt

import (
	"net/http"

	"github.com/severedsea/golang-kit/web"
)

var (
	// 4xx

	// ErrInvalidToken is the error returned if the access token string is invalid
	ErrInvalidToken = &web.Error{Status: http.StatusUnauthorized, Code: "invalid_token", Desc: "Invalid access token"}

	// 5xx

	// ErrJWT is the error returned for any unexpected error related to generation of JWT
	ErrJWT = &web.Error{Status: http.StatusInternalServerError, Code: "jwt"}
)
