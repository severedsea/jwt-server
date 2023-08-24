package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v4"
	"github.com/severedsea/golang-kit/ptr"
	"github.com/severedsea/golang-kit/web"
	"github.com/severedsea/jwt-server/internal/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const tokenString = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhYyIsInN1YiI6ImVtYWlsIn0.ABKz3XHisql1kJOS59K_iYCT5ltjfo_I0zlohc8Nuhioi_4sIXiUXJyIQJ0YCfeYBuLq8BnNmzagNbvwdyAWf2A1jRfePG6Gtn216tlbeH-EF-FD7Z4Z6fXPmXWNOEbTR9wlPBS2WyJOz1y2jJLTYKn70VYwIEqDY0UAIY8QyDAdds-TMA27XnvP6gkVcAbbu-kE35q5C3bAK-ZyjULlm-LL0JPJIAyS_3RxHZ1xKxLHH16iVAkdUsbi1fNVCi44yGCIglrt77gvUBpvRY12D8mERlc5DqbgdKI8W1FiNGMrSXsaBM7mt7N8XGYHA4OnWQl6giyKoKzFyeG9SPTJKg"

func TestMiddleware(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc       string
		cookie     *http.Cookie
		authHeader *string
	}{
		{
			desc:       "cookie only",
			cookie:     &http.Cookie{Name: tokenCookieName, Value: tokenString},
			authHeader: nil,
		},
		{
			desc:       "header only",
			cookie:     nil,
			authHeader: ptr.Reference("Bearer " + tokenString),
		},
		{
			desc:       "cookie + header",
			cookie:     &http.Cookie{Name: tokenCookieName, Value: tokenString},
			authHeader: ptr.Reference("Bearer " + tokenString),
		},
		{
			desc:       "INVALID cookie + header",
			cookie:     &http.Cookie{Name: tokenCookieName, Value: "INVALID"},
			authHeader: ptr.Reference("Bearer " + tokenString),
		},
		{
			desc:       "cookie + INVALID header",
			cookie:     &http.Cookie{Name: tokenCookieName, Value: tokenString},
			authHeader: ptr.Reference("INVALID"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			var passed bool
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				passed = true
				ctx := r.Context()
				_, err := ClaimsFromContext(ctx)
				assert.NoError(t, err)
			})
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/some/path", nil)

			if tc.authHeader != nil {
				r.Header.Set("Authorization", *tc.authHeader)
			}

			if tc.cookie != nil {
				r.AddCookie(tc.cookie)
			}

			// Mocks:
			stub := &mockTokenParserVerifier{}
			stub.On("ParseToken", mock.Anything, tokenString).
				Return(Claims{
					RegisteredClaims: jwtgo.RegisteredClaims{Subject: "SUBJECT"},
				}, nil)
			stub.On("VerifyToken", mock.Anything, tokenString, "SUBJECT").
				Return(nil)

			// When:
			Middleware(stub)(handler).ServeHTTP(w, r)

			// Then:
			assert.True(t, passed)
			stub.AssertCalled(t, "ParseToken", mock.Anything, tokenString)
		})
	}
}

func TestMiddleware_Error(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc        string
		cookie      *http.Cookie
		authHeader  *string
		parseCalls  int
		parseErr    error
		verifyCalls int
		verifyErr   error
		expected    *web.Error
	}{
		{
			desc:       "empty cookie only",
			cookie:     &http.Cookie{Name: tokenCookieName, Value: ""},
			authHeader: nil,
			expected:   ErrMissingToken,
		},
		{
			desc:       "empty header only",
			cookie:     nil,
			authHeader: ptr.Reference("Bearer"),
			expected:   ErrMissingToken,
		},
		{
			desc:       "no cookie + no header",
			cookie:     nil,
			authHeader: nil,
			expected:   ErrMissingToken,
		},
		{
			desc:       "empty cookie + empty header",
			cookie:     &http.Cookie{Name: tokenCookieName, Value: ""},
			authHeader: ptr.Reference(""),
			expected:   ErrMissingToken,
		},
		{
			desc:       "parse token failed",
			cookie:     &http.Cookie{Name: tokenCookieName, Value: ""},
			authHeader: ptr.Reference("Bearer " + tokenString),
			parseErr:   jwt.ErrInvalidToken,
			parseCalls: 1,
			expected:   jwt.ErrInvalidToken,
		},
		{
			desc:        "verify token failed",
			cookie:      &http.Cookie{Name: tokenCookieName, Value: "TOKEN"},
			authHeader:  ptr.Reference("Bearer " + tokenString),
			verifyErr:   jwt.ErrInvalidToken,
			verifyCalls: 1,
			parseCalls:  1,
			expected:    jwt.ErrInvalidToken,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Fail()
			})
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/some/path", nil)

			if tc.authHeader != nil {
				r.Header.Set("Authorization", *tc.authHeader)
			}

			if tc.cookie != nil {
				r.AddCookie(tc.cookie)
			}

			// Mocks:
			stub := &mockTokenParserVerifier{}
			stub.On("ParseToken", mock.Anything, mock.Anything).
				Return(Claims{}, tc.parseErr)
			stub.On("VerifyToken", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.verifyErr)

			// When:
			Middleware(stub)(handler).ServeHTTP(w, r)

			// Then:
			stub.AssertNumberOfCalls(t, "ParseToken", tc.parseCalls)
			stub.AssertNumberOfCalls(t, "VerifyToken", tc.verifyCalls)

			b, err := io.ReadAll(w.Result().Body)
			assert.NoError(t, err)
			var actual web.Error
			assert.NoError(t, json.Unmarshal(b, &actual))
			assert.Equal(t, tc.expected.Code, actual.Code)
			assert.Equal(t, tc.expected.Desc, actual.Desc)
			assert.Equal(t, tc.expected.Status, w.Result().StatusCode)

			// Assert invalidated cookie
			var found bool
			for _, it := range w.Result().Cookies() {
				if it.Name == tokenCookieName {
					assert.Equal(t, true, it.HttpOnly)
					assert.True(t, it.Expires.Before(time.Now()))
					assert.Equal(t, "deleted", it.Value)
					found = true
				}
			}
			assert.True(t, found)
		})
	}
}

// mockTokenParserVerifier is the mock token parser
type mockTokenParserVerifier struct {
	mock.Mock
}

func (m *mockTokenParserVerifier) ParseToken(ctx context.Context, tokenString string) (Claims, error) {
	args := m.Called(ctx, tokenString)

	return args.Get(0).(Claims), args.Error(1)
}

func (m *mockTokenParserVerifier) VerifyToken(ctx context.Context, tokenString, subject string) error {
	args := m.Called(ctx, tokenString, subject)

	return args.Error(0)
}
