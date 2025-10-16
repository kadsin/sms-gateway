package userbalance

import (
	"context"

	"github.com/google/uuid"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/redis/go-redis/v9"
)

func Get(ctx context.Context, userId uuid.UUID) (float32, error) {
	rKey := CacheKey(userId)

	tx := container.Redis().Get(ctx, rKey)
	if tx.Err() == redis.Nil {
		balance, err := getRealUserBalance(ctx, userId)
		if err != nil {
			return 0, err
		}

		container.Redis().Set(ctx, rKey, balance, 0)
		return balance, nil
	}

	f, err := tx.Float32()
	if err != nil {
		return 0, err
	}

	return f, nil
}
