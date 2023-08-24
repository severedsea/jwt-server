package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	// Given:

	// When:
	c, err := New()

	// Then:
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestRedisURLFromEnv(t *testing.T) {
	testCases := []struct {
		scheme   string
		user     string
		password string
		host     string
		port     string
		expected string
	}{
		{
			scheme:   "rediss",
			user:     "user",
			password: "password",
			host:     "localhost",
			port:     "1433",
			expected: "rediss://user:password@localhost:1433", // pragma: allowlist secret
		},
		{
			scheme:   "rediss",
			user:     "user",
			password: "",
			host:     "localhost",
			port:     "1433",
			expected: "rediss://user:@localhost:1433",
		},
		{
			scheme:   "rediss",
			user:     "",
			password: "password",
			host:     "localhost",
			port:     "1433",
			expected: "rediss://:password@localhost:1433",
		},
		{
			scheme:   "rediss",
			user:     "",
			password: "",
			host:     "localhost",
			port:     "1433",
			expected: "rediss://localhost:1433",
		},
		{
			scheme:   "rediss",
			user:     "",
			password: "",
			host:     "localhost",
			port:     "",
			expected: "rediss://localhost",
		},
		{
			scheme:   "",
			user:     "",
			password: "",
			host:     "localhost",
			port:     "1433",
			expected: "localhost:1433",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			// Given:
			t.Setenv("REDIS_SCHEME", tc.scheme)
			t.Setenv("REDIS_USER", tc.user)
			t.Setenv("REDIS_PWD", tc.password)
			t.Setenv("REDIS_HOST", tc.host)
			t.Setenv("REDIS_PORT", tc.port)

			// When:
			actual := redisURLFromEnv()

			// Then:
			assert.Equal(t, tc.expected, actual)
		})
	}
}
