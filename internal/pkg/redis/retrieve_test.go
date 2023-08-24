package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieveKeysByPattern(t *testing.T) {
	ctx := context.Background()

	redisClient, err := New()
	require.NoError(t, err)

	// Given:
	given := []string{
		"mtiuen_03100001X",
		"mtiuen_03100002X",
		"mtiuen_03100003X",
		"mtiuen_03100004X",
		"mtiuen_03100005X",
		"mtiuen_03100006X",
		"mtiuen_03100007X",
		"mtiuen_03100008X",
	}
	for _, it := range given {
		err = redisClient.Set(ctx, it, "", 0).Err()
		require.NoError(t, err)
	}

	// When:
	act, err := RetrieveKeysByPattern(ctx, redisClient, "mtiuen_*", 5)
	require.NoError(t, err)

	// Then:
	assert.Equal(t, 8, len(act))
	assert.ElementsMatch(t, given, act)
}
