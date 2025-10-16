package userbalance

import (
	"context"

	"github.com/google/uuid"
	"github.com/kadsin/sms-gateway/internal/container"
)

func Change(ctx context.Context, userID uuid.UUID, delta float32) (float32, error) {
	rKey := CacheKey(userID)

	exists, err := container.Redis().Exists(ctx, rKey).Result()
	if err != nil {
		return 0, err
	}

	if exists == 0 {
		if _, err := syncRedisByPostgre(ctx, userID); err != nil {
			return 0, err
		}
	}

	newBalance, err := container.Redis().IncrByFloat(ctx, rKey, float64(delta)).Result()
	if err != nil {
		return 0, err
	}

	return float32(newBalance), nil
}
