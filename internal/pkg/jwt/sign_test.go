package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/severedsea/golang-kit/web"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	t.Parallel()

	// Given:
	given := jwt.RegisteredClaims{
		Subject: "subject",
		Issuer:  Issuer,
	}

	// When:
	actual, err := Sign(given)

	// Then:
	assert.NoError(t, err)
	assert.NotEmpty(t, actual)

	// Verify access token
	result := jwt.RegisteredClaims{}
	assert.NoError(t, Parse(actual, &result))
	assert.Equal(t, given, result)
}

func TestParse(t *testing.T) {
	testCases := []struct {
		desc     string
		given    jwt.RegisteredClaims
		expected *web.Error
	}{
		{
			desc: "no issuer",
			given: jwt.RegisteredClaims{
				Subject: "subject",
			},
			expected: ErrInvalidToken,
		},
		{
			desc: "expired",
			given: jwt.RegisteredClaims{
				Subject:   "subject",
				Issuer:    Issuer,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Second)),
			},
			expected: ErrInvalidToken,
		},
		{
			desc: "future issued at",
			given: jwt.RegisteredClaims{
				Subject:   "subject",
				Issuer:    Issuer,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
			expected: ErrInvalidToken,
		},
		{
			desc: "valid",
			given: jwt.RegisteredClaims{
				Subject:   "subject",
				Issuer:    Issuer,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			given, err := Sign(tc.given)
			assert.NoError(t, err)

			// When:
			c := jwt.RegisteredClaims{}
			err = Parse(given, &c)

			// Then:
			if err != nil {
				assert.Equal(t, tc.expected, err)
			} else {
				assert.Equal(t, tc.given, c)
			}
		})
	}
}

func TestParse_Invalid_Error(t *testing.T) {
	// Given:

	// When:
	c := jwt.RegisteredClaims{}
	err := Parse("INVALID", &c)

	// Then:
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidToken, err)
}
