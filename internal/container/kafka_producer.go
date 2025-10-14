package container

import (
	"time"

	"github.com/kadsin/sms-gateway/config"
	"github.com/kadsin/sms-gateway/internal/qkafka"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

var kafkaProducerInstance qkafka.Producer

func initKafkaProducer() {
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

	kafkaProducerInstance = qkafka.NewProducer(kafka.WriterConfig{
		Brokers: config.Env.Kafka.Brokers,

		Balancer: &kafka.LeastBytes{},
		Dialer:   dialer,
		Async:    false,
	})
}

func KafkaProducer() qkafka.Producer {
	if dbInstance == nil {
		fatalOnNilInstance("Kafka producer is not connected.")
	}

	return kafkaProducerInstance
}
