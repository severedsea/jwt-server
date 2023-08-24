package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestClaimsContext(t *testing.T) {
	t.Parallel()

	// Given:
	r := httptest.NewRequest(http.MethodGet, "/some/path", nil)
	given := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "christina_ang",
		},
	}

	// When: Set
	newCtx := setClaimsContext(r.Context(), given)

	// When: Get
	result, err := ClaimsFromContext(newCtx)

	// Then:
	assert.NoError(t, err)
	assert.Equal(t, given, result)
}

func TestClaimsContext_Empty(t *testing.T) {
	t.Parallel()

	// Given:
	r := httptest.NewRequest(http.MethodGet, "/some/path", nil)
	ctx := r.Context()

	// When: Get
	_, err := ClaimsFromContext(ctx)

	// Then:
	assert.Equal(t, ErrMissingContext, err)
}
