package middlewares

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kadsin/sms-gateway/config"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/segmentio/kafka-go"
	"github.com/sony/gobreaker"
)

func NewCircuitBreaker() *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        "SystemHealthChecker",
		MaxRequests: 3,
		Interval:    0,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 10
		},
	}

	return gobreaker.NewCircuitBreaker(settings)
}

func CircuitBreakerMiddleware(cb *gobreaker.CircuitBreaker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := cb.Execute(func() (any, error) {
			return nil, checkSystemHealth(c.Context())
		})

		if err != nil {
			return err
		}

		return c.Next()
	}
}

func checkSystemHealth(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	if err := container.Redis().Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping error: %w", err)
	}

	conn, err := kafka.DialContext(ctx, "tcp", config.Env.Kafka.Brokers[0])
	if err != nil {
		return fmt.Errorf("kafka broker not reachable: %w", err)
	}
	defer conn.Close()

	return nil
}
