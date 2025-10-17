package wallet

import (
	"context"

	"github.com/google/uuid"
	"github.com/kadsin/sms-gateway/config"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
	"github.com/segmentio/kafka-go"
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

	publishOnKafka(ctx, messages.WalletChanged{
		ClientId: userID,
		Amount:   delta,
	})

	return float32(newBalance), nil
}

func publishOnKafka(ctx context.Context, ub messages.WalletChanged) error {
	b, err := dtos.Marshal(ub)
	if err != nil {
		return err
	}

	return container.KafkaProducer().SendMessage(ctx, kafka.Message{
		Topic: config.Env.Kafka.Topics.Balance,
		Value: b,
	})
}
