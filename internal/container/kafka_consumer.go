package container

import (
	"time"

	"github.com/kadsin/sms-gateway/config"
	"github.com/kadsin/sms-gateway/internal/qkafka"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func KafkaConsumer(topic string) qkafka.Consumer {
	if mockKafkaConsumerInstance != nil {
		return mockKafkaConsumerInstance
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	if config.Env.Kafka.Username != "" && config.Env.Kafka.Password != "" {
		dialer.SASLMechanism = plain.Mechanism{
			Username: config.Env.Kafka.Username,
			Password: config.Env.Kafka.Password,
		}
	}

	return qkafka.NewConsumer(kafka.ReaderConfig{
		Brokers: config.Env.Kafka.Brokers,
		Topic:   topic,
		GroupID: "consumergroup." + topic,
		Dialer:  dialer,
	})
}
