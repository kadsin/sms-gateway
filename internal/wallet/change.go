package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/google/uuid"
	"github.com/kadsin/sms-gateway/config"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
	"github.com/segmentio/kafka-go"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrLockNotObtained   = errors.New("could not obtain wallet lock")
)

func Change(ctx context.Context, userID uuid.UUID, delta float32) (float32, error) {
	rKey := CacheKey(userID)
	rdb := container.Redis()

	lock, err := redislock.New(rdb).Obtain(
		ctx,
		fmt.Sprintf("lock:wallet:%s", userID.String()),
		5*time.Second,
		&redislock.Options{RetryStrategy: redislock.LinearBackoff(100 * time.Millisecond)},
	)
	if err != nil {
		return 0, ErrLockNotObtained
	}
	defer lock.Release(ctx)

	oldBalance, err := Get(ctx, userID)
	if err != nil {
		return 0, err
	}

	newBalance := oldBalance + delta
	if newBalance < 0 {
		return oldBalance, ErrInsufficientFunds
	}

	if err := rdb.IncrByFloat(ctx, rKey, float64(delta)).Err(); err != nil {
		return oldBalance, err
	}

	err = publishOnKafka(ctx, messages.WalletChanged{
		ClientId: userID,
		Amount:   delta,
	})
	if err != nil {
		return oldBalance, err
	}

	return newBalance, nil
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
