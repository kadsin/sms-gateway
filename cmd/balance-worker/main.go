package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kadsin/sms-gateway/config"
	"github.com/kadsin/sms-gateway/database/models"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
	"github.com/kadsin/sms-gateway/internal/qkafka"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	container.Init()
	defer container.Close()

	if len(os.Args) < 2 {
		log.Fatal("Consumers count is required as the first argument")
	}
	arg1 := strings.Trim(os.Args[1], " \t\n\r\"'")
	consumers, err := strconv.Atoi(arg1)
	if err != nil {
		log.Fatalf("Error on using argument with value %v as int for consumers count", arg1)
	}

	log.Printf("Initializing consumer pool with count %v ...", consumers)
	g, ctx := errgroup.WithContext(context.TODO())

	for i := range consumers {
		cid := i
		consumer := container.KafkaConsumer(config.Env.Kafka.Topics.Balance)

		g.Go(func() error {
			defer consumer.Close()
			log.Printf("consumer %d listening on topic `%s` with group `%s` ...", cid, consumer.Config().Topic, consumer.Config().GroupID)

			err := processMessages(ctx, cid, consumer)
			if err != nil {
				log.Printf("consumer %d stopped with error: %v", cid, err)
			}
			return err
		})
	}

	if err := g.Wait(); err != nil {
		log.Fatalf("App stopped because of consumer error: %v", err)
	}
}

func processMessages(ctx context.Context, consumerId int, consumer qkafka.Consumer) error {
	for {
		m, err := consumer.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}

			catchKafkaError(err)
			continue
		}

		ub, err := dtos.Unmarshal[messages.UserBalance](m.Value)
		if err != nil {
			log.Printf("Error on unmarshalling message: %v", err)
			continue
		}

		if err := updateBalanceInPostgres(ctx, ub); err != nil {
			return fmt.Errorf("consumer id: %v\nerror:%+v", consumerId, err)
		}

		log.Printf("consumer id: %v\nbalance updated\nid: %s, amount: %v\n", consumerId, ub.ClientId, ub.Amount)

		if err := consumer.Commit(ctx, m); err != nil {
			return fmt.Errorf("error on committing message: %v", err)
		}
	}
}

func catchKafkaError(err error) {
	// TODO: Should capture by sentry
	log.Printf("Error on fetching message from kafka:\n%+v", err)

	log.Println("Sleep 2 seconds ...")
	time.Sleep(2 * time.Second)
	log.Println("Fetching ...")
}

func updateBalanceInPostgres(ctx context.Context, ub messages.UserBalance) error {
	return container.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&models.User{}).
			Where("id", ub.ClientId).
			UpdateColumn("balance", gorm.Expr("COALESCE(balance, 0) + (?)", ub.Amount)).
			Error
	})
}
