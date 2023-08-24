package auth

import (
	"net/http"
	"os"
)

const (
	// TokenCookieName is the recommended token cookie name
	tokenCookieName = "token"
)

func newCookie() http.Cookie {
	cookie := http.Cookie{
		Name:     tokenCookieName,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	if os.Getenv("APP_ENV") == "local" {
		cookie.Secure = false
		cookie.Domain = "localhost"
		cookie.SameSite = http.SameSiteLaxMode
	}

	return cookie
}

// Cookie returns an HTTP cookie from the Token values
func (t Token) Cookie() *http.Cookie {
	cookie := newCookie()
	cookie.Value = t.AccessToken
	cookie.Expires = t.ExpiresAt

	return &cookie
}

// InvalidateCookie returns a cookie meant to invalidate the auth cookie
func InvalidateCookie(w http.ResponseWriter) {
	cookie := newCookie()
	cookie.Value = "deleted"
	cookie.MaxAge = -1

	// Invalidate token cookie for web clients
	http.SetCookie(w, &cookie)
}
