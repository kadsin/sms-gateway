package qkafka

import (
	"context"
	"time"

	"github.com/kadsin/sms-gateway/config"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

var producerInstance producer = nil

func init() {
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

	producerInstance = &KafkaProducer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: config.Env.Kafka.Brokers,

			Balancer: &kafka.LeastBytes{},
			Dialer:   dialer,
			Async:    false,
		}),
	}
}

func Producer() producer {
	return producerInstance
}

type producer interface {
	SendMessage(ctx context.Context, m kafka.Message) error
	Close() error
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func (p *KafkaProducer) SendMessage(ctx context.Context, m kafka.Message) error {
	return p.writer.WriteMessages(ctx, m)
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
