package container_test

import (
	"context"
	"os"
	"testing"

	"github.com/kadsin/sms-gateway/internal/container"
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

func TestMain(m *testing.M) {
	container.Init()
	defer container.Close()

	m.Run()

	os.Exit(0)
}

func Test_Mock_KafkaProducer(t *testing.T) {
	container.MockKafkaProducer(&ProducerMock{})

	p := container.KafkaProducer()
	require.IsType(t, &ProducerMock{}, p)
}

func Test_Mock_KafkaConsumer(t *testing.T) {
	container.MockKafkaConsumer(&ConsumerMock{})

	p := container.KafkaConsumer("test")
	require.IsType(t, &ConsumerMock{}, p)
}
