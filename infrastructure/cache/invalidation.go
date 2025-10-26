package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// InvalidateOnboardingStatusCache invalidates the onboarding-status cache for a specific user
func InvalidateOnboardingStatusCache(ctx context.Context, rdb *redis.Client, userID string) error {
	key := OnboardingStatusCacheKey(userID)
	_, err := Delete(ctx, rdb, key)
	return err
}
