package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/severedsea/golang-kit/web"
	"github.com/severedsea/golang-kit/web/middleware"
)

// Middleware parses the bearer Authorization or cookie, and validates the JWT signature
func Middleware(p TokenParserVerifier) middleware.Adapter {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			c, err := validateToken(ctx, r, p)
			if err != nil {
				InvalidateCookie(w)
				web.RespondJSON(ctx, w, err, nil)

				return
			}

			ctx = setClaimsContext(ctx, c)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func validateToken(ctx context.Context, r *http.Request, p TokenParserVerifier) (Claims, error) {
	token := tokenFromRequest(r)
	if token == "" {
		return Claims{}, ErrMissingToken
	}

	c, err := p.ParseToken(ctx, token)
	if err != nil {
		return Claims{}, err
	}

	if err := p.VerifyToken(ctx, token, c.Subject); err != nil {
		return Claims{}, err
	}

	return c, nil
}

// tokenFromRequest tries to retrieve the token string by calling the token funcs in order
func tokenFromRequest(r *http.Request) string {
	var tokenString string

	findToken := []func(r *http.Request) string{
		tokenFromHeader, tokenFromCookie,
	}

	for _, f := range findToken {
		tokenString = f(r)
		if tokenString != "" {
			break
		}
	}

	return strings.TrimSpace(tokenString)
}

// tokenFromCookie tries to retrieve the token string from a cookie named
// "token".
func tokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(tokenCookieName)
	if err != nil {
		return ""
	}

	return cookie.Value
}

// tokenFromHeader tries to retrieve the token string from the
// "Authorization" request header: "Authorization: BEARER T".
func tokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}

	return ""
}
