package container

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kadsin/sms-gateway/config"
	"github.com/redis/go-redis/v9"
)

var redisInstance *redis.Client

func initRedis() {
	addr := fmt.Sprintf("%s:%s", config.Env.Redis.Host, config.Env.Redis.Port)

	redisInstance = redis.NewClient(&redis.Options{
		Addr:         addr,
		Username:     config.Env.Redis.Username,
		Password:     config.Env.Redis.Password,
		DB:           config.Env.Redis.Database,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     50,
		MinIdleConns: 10,
	})

	if err := redisInstance.Ping(context.TODO()).Err(); err != nil {
		log.Fatalf("Redis connection failed at %s: %v", addr, err)
	}
}

func Redis() *redis.Client {
	if redisInstance == nil {
		fatalOnNilInstance("Redis is not connected.")
	}

	return redisInstance
}
