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

func main() {
	container.Init()
	defer container.Close()

	processMessages(context.TODO())
}

func processMessages(ctx context.Context) error {
	topic := strings.Trim(os.Args[1], " \t\n\r\"'")
	log.Printf("Start listening on topic %s ...", topic)

	consumer := container.KafkaConsumer(topic)
	smsProvider := container.SmsProvider()

	for {
		m, err := consumer.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return err
			}

			catchKafkaError(err)
			continue
		}

		s, err := dtos.Unmarshal[messages.Sms](m.Value)
		if err != nil {
			log.Printf("Error on unmarhsaling kafka message with value: \n%s\nerror:%+v", m.Value, err)
		}
		go sendSms(smsProvider, s)
	}

}

func catchKafkaError(err error) {
	// TODO: Should capture by sentry
	log.Printf("Error on fetching message from kafka:\n%+v", err)

	log.Println("Sleep 2 seconds ...")
	time.Sleep(2 * time.Second)
	log.Println("Fetching ...")
}

func sendSms(p sms.Provider, m messages.Sms) {
	if err := p.Send(m); err != nil {
		log.Printf("Error on send sms with body: %+v\nerror:%+v", m, err)

		analytics_models.LogFailure(&m)
		return
	}

	analytics_models.LogSent(&m)
}
