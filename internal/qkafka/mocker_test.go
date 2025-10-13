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

func Test_MockQKafkaProducer(t *testing.T) {
	qkafka.MockProducer(&ProducerMock{})

	p := qkafka.NewProducer(kafka.WriterConfig{})
	require.IsType(t, &ProducerMock{}, p)
}
