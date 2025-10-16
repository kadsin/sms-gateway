package userbalance

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kadsin/sms-gateway/database/models"
	"github.com/kadsin/sms-gateway/internal/container"
	"gorm.io/gorm/clause"
)

var group = &lockGroup{locks: make(map[uuid.UUID]*userLock)}

func init() {
	go func() {
		for {
			time.Sleep(cleanupInterval)

			group.cleanup()
		}
	}()
}

func UserLock(id uuid.UUID) *sync.Mutex {
	return group.GetOrCreate(id)
}

func CacheKey(userID uuid.UUID) string {
	return fmt.Sprintf("user_balance:%s", userID)
}

func syncRedisByPostgre(ctx context.Context, userId uuid.UUID) (float32, error) {
	balance, err := getRealUserBalance(ctx, userId)
	if err != nil {
		return 0, err
	}

	container.Redis().Set(ctx, CacheKey(userId), balance, 0)
	return balance, nil
}

func getRealUserBalance(_ context.Context, userId uuid.UUID) (float32, error) {
	tx := container.DB().Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer tx.Rollback()

	var user models.User
	err := tx.Clauses(clause.Locking{Strength: "SHARE"}).
		Where("id", userId).
		First(&user).Error
	if err != nil {
		return 0, err
	}

	tx.Commit()

	return user.Balance, nil
}
