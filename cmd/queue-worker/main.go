package main

import (
	"context"
	"errors"
	"log"
	"math"
	"os"
	"strings"
	"time"

	analytics_models "github.com/kadsin/sms-gateway/analytics/models"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
	"github.com/kadsin/sms-gateway/internal/qkafka"
	"github.com/kadsin/sms-gateway/internal/sms"
	userbalance "github.com/kadsin/sms-gateway/internal/user_balance"
)

const WORKER_COUNT = 50

func main() {
	container.Init()
	defer container.Close()

	if len(os.Args) < 2 {
		log.Fatal("Topic name is required as the first argument")
	}
	topic := strings.Trim(os.Args[1], " \t\n\r\"'")

	consumer := container.KafkaConsumer(topic)
	defer consumer.Close()
	log.Printf("Start listening on topic `%s` with group id `%s` ...", consumer.Config().Topic, consumer.Config().GroupID)

	processMessages(context.TODO(), consumer)
}

func processMessages(ctx context.Context, consumer qkafka.Consumer) error {
	log.Printf("Initializing workers pool with count %v ...", WORKER_COUNT)
	msgChan := initialSmsWorkersPool(ctx)

	for {
		m, err := consumer.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				close(msgChan)
				return err
			}

			catchKafkaError(err)
			continue
		}

		if err := consumer.Commit(ctx, m); err != nil {
			log.Printf("Error on committing message: %v", err)
		}

		s, err := dtos.Unmarshal[messages.Sms](m.Value)
		if err != nil {
			log.Printf("Error on unmarshalling message: %v", err)
			continue
		}

		msgChan <- s
	}
}

func initialSmsWorkersPool(ctx context.Context) chan messages.Sms {
	provider := container.SmsProvider()

	msgChan := make(chan messages.Sms, 1000)

	for i := range WORKER_COUNT {
		go func(id int) {
			for m := range msgChan {
				if err := sendSms(ctx, provider, m); err != nil {
					log.Printf("worker id: %v\nerror:%+v", id, err)
				} else {
					log.Printf("worker id: %v\nsms sent\nid: %s\n", id, m.Id)
				}
			}
		}(i)
	}

	return msgChan
}

func sendSms(ctx context.Context, p sms.Provider, m messages.Sms) error {
	if err := p.Send(m); err != nil {
		analytics_models.LogFailure(&m)

		amount := math.Abs(float64(m.Price))
		if _, err := userbalance.Change(ctx, m.SenderClientId, float32(amount)); err != nil {
			return err
		}

		return err
	}

	analytics_models.LogSent(&m)
	return nil
}

func catchKafkaError(err error) {
	// TODO: Should capture by sentry
	log.Printf("Error on fetching message from kafka:\n%+v", err)

	log.Println("Sleep 2 seconds ...")
	time.Sleep(2 * time.Second)
	log.Println("Fetching ...")
}
