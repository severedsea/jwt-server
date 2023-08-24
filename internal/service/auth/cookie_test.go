package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToken_Cookie(t *testing.T) {
	t.Parallel()

	// Given:
	given := Token{
		AccessToken: "access_token",
		ExpiresIn:   1,
		ExpiresAt:   time.Now().Add(time.Second),
		Scope:       "some_scope",
		TokenType:   "Bearer",
	}

	// When:
	actual := given.Cookie()

	// Then:
	assert.NotNil(t, actual)
	assert.Equal(t, true, actual.HttpOnly)
	assert.Equal(t, true, actual.Secure)
	assert.Equal(t, http.SameSiteStrictMode, actual.SameSite)
	assert.Equal(t, given.ExpiresAt.Unix(), actual.Expires.Unix())
	assert.Equal(t, given.AccessToken, actual.Value)
	assert.Equal(t, "", actual.Domain)
	assert.Equal(t, "/", actual.Path)
	assert.Equal(t, 0, actual.MaxAge)
}
