package userbalance

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kadsin/sms-gateway/database/models"
	"github.com/kadsin/sms-gateway/internal/container"
	"gorm.io/gorm/clause"
)

func CacheKey(userID uuid.UUID) string {
	return fmt.Sprintf("user_balance:%s", userID)
}

func getRealUserBalance(ctx context.Context, userId uuid.UUID) (float32, error) {
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
