package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client with the provided address and password.
// addr should be in the form "host:port" (e.g. "localhost:6379").
func NewRedisClient(addr, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
}

// SetWithTTL sets a key with a TTL. Returns any error from Redis.
func SetWithTTL(ctx context.Context, rdb *redis.Client, key string, value string, ttl time.Duration) error {
	return rdb.Set(ctx, key, value, ttl).Err()
}

// Get gets a value by key. If key does not exist, returns redis.Nil.
func Get(ctx context.Context, rdb *redis.Client, key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

// Delete removes the given key. Returns number of removed keys and error.
func Delete(ctx context.Context, rdb *redis.Client, key string) (int64, error) {
	return rdb.Del(ctx, key).Result()
}

// OnboardingStatusCacheKey returns the Redis key used to store onboarding-status for a user.
func OnboardingStatusCacheKey(userID string) string {
	return "onboarding_status:user:" + userID
}
