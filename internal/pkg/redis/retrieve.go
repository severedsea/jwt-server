package redis

import (
	"context"
	"sort"

	"github.com/go-redis/redis/v8"
	"golang.org/x/exp/slices"
)

// RetrieveKeysByPattern returns all keys that match the specified pattern from redis
func RetrieveKeysByPattern(ctx context.Context, r redis.Cmdable, pattern string, batchSize int64) ([]string, error) {
	var keys, allKeys []string
	var cursor uint64
	var err error

	for {
		keys, cursor, err = r.Scan(ctx, cursor, pattern, batchSize).Result()
		if err != nil {
			return nil, err
		}

		// Scan may return no keys even though the iteration has not completed
		// https://redis.io/commands/scan/
		if len(keys) > 0 {
			allKeys = append(allKeys, keys...)
		}

		if cursor == 0 {
			// Dedup as Scan may return duplicates
			// https://redis.io/commands/scan/
			sort.Strings(allKeys)
			return slices.Compact(allKeys), nil
		}
	}
}
