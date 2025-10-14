package container_test

import (
	"os"
	"testing"

	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/tests/mocks"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	container.Init()
	defer container.Close()

	m.Run()

	os.Exit(0)
}

func Test_Mock_KafkaProducer(t *testing.T) {
	container.MockKafkaProducer(&mocks.KafkaProducerMock{})

	p := container.KafkaProducer()
	require.IsType(t, &mocks.KafkaProducerMock{}, p)
}

func Test_Mock_KafkaConsumer(t *testing.T) {
	container.MockKafkaConsumer(&mocks.KafkaConsumerMock{})

	p := container.KafkaConsumer("test")
	require.IsType(t, &mocks.KafkaConsumerMock{}, p)
}
