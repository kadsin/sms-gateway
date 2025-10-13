package qkafka_test

import (
	"context"
	"testing"

	"github.com/kadsin/sms-gateway/internal/qkafka"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

type ProducerMock struct{}

func (mp *ProducerMock) SendMessage(ctx context.Context, m kafka.Message) error {
	return nil
}

func (mp *ProducerMock) Close() error {
	return nil
}

type ConsumerMock struct{}

func (mp *ConsumerMock) FetchMessage(ctx context.Context) (kafka.Message, error) {
	return kafka.Message{}, nil
}

func (mp *ConsumerMock) Commit(ctx context.Context, msgs ...kafka.Message) error {
	return nil
}

func (mp *ConsumerMock) Close() error {
	return nil
}

func Test_MockQKafkaProducer(t *testing.T) {
	qkafka.MockProducer(&ProducerMock{})

	p := qkafka.NewProducer(kafka.WriterConfig{})
	require.IsType(t, &ProducerMock{}, p)
}

func Test_MockQKafkaConsumer(t *testing.T) {
	qkafka.MockConsumer(&ConsumerMock{})

	p := qkafka.NewConsumer(kafka.ReaderConfig{})
	require.IsType(t, &ConsumerMock{}, p)
}
