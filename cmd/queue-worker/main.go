package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	analytics_models "github.com/kadsin/sms-gateway/analytics/models"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
	"github.com/kadsin/sms-gateway/internal/sms"
)

const WORKER_COUNT = 50

func main() {
	container.Init()
	defer container.Close()

	if len(os.Args) < 2 {
		log.Fatal("Topic name is required as the first argument")
	}
	topic := strings.Trim(os.Args[1], " \t\n\r\"'")

	processMessages(context.TODO(), topic)
}

func processMessages(ctx context.Context, topic string) error {
	log.Printf("Initializing workers pool with count %v ...", WORKER_COUNT)
	msgChan := initialSmsWorkersPool()

	log.Printf("Start listening on topic `%s` ...", topic)
	consumer := container.KafkaConsumer(topic)
	defer consumer.Close()

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

func initialSmsWorkersPool() chan messages.Sms {
	provider := container.SmsProvider()

	msgChan := make(chan messages.Sms, 1000)

	for i := range WORKER_COUNT {
		go func(id int) {
			for m := range msgChan {
				if err := sendSms(provider, m); err != nil {
					log.Printf("worker id: %v\nerror:%+v", id, err)
				} else {
					log.Printf("worker id: %v\nsms sent\nid: %s\n", id, m.Id)
				}
			}
		}(i)
	}

	return msgChan
}

func sendSms(p sms.Provider, m messages.Sms) error {
	if err := p.Send(m); err != nil {
		analytics_models.LogFailure(&m)
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
