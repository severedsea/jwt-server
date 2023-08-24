package auth

import (
	"net/http"

	"github.com/severedsea/golang-kit/web"
)

var (
	// ErrMissingContext is the error returned if there's no auth in the context
	// This is mainly due to implementation error
	ErrMissingContext = &web.Error{Status: http.StatusInternalServerError, Code: "missing_auth_context", Desc: "Missing auth context"}
	// ErrMissingToken is the error returned if the access token string is missing
	ErrMissingToken = &web.Error{Status: http.StatusUnauthorized, Code: "missing_token", Desc: "Missing access token"}
	// ErrInactiveToken is the error returned if the token retrieved using the authorization code is inactive
	ErrInactiveToken = &web.Error{Status: http.StatusBadRequest, Code: "inactive_token", Desc: "inactive token"}
	// ErrRedis is the generic web error for redis-related errors
	ErrRedis = &web.Error{Status: http.StatusInternalServerError, Code: "redis"}
	// ErrInternal is the generic web error for internal errors
	ErrInternal = &web.Error{Status: http.StatusInternalServerError, Code: "internal"}
)
