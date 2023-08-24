package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRegisteredClaims(t *testing.T) {
	// Given:
	username := "christina_ang"
	duration := time.Hour

	// When:
	actual := NewRegisteredClaims(username, duration)

	// Then:
	afterExecTime := time.Now()
	assert.Equal(t, username, actual.Subject)
	assert.Equal(t, Issuer, actual.Issuer)
	assert.True(t, actual.IssuedAt.Unix() > 0, "should be populated")
	assert.True(t, actual.IssuedAt.Unix() > time.Time{}.Unix(), "should not be zero time")
	assert.True(t, actual.IssuedAt.Unix() >= afterExecTime.Unix())
	assert.True(t, actual.ExpiresAt.Unix() > 0, "should be populated")
	assert.True(t, actual.ExpiresAt.Unix() > time.Time{}.Unix(), "should not be zero time")
	assert.True(t, actual.ExpiresAt.Unix() <= afterExecTime.Add(duration).Unix())
}
