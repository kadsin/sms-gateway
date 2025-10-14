package tests

import (
	"os"
	"testing"

	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/tests/mocks"
)

func TestMain(m *testing.M) {
	container.Init()
	defer container.Close()

	mockEssentials()

	RefreshDatabase()

	container.DB().Begin()
	code := m.Run()
	container.DB().Rollback()

	os.Exit(code)
}

func mockEssentials() {
	container.MockKafkaProducer(&mocks.KafkaProducerMock{})

	container.MockKafkaConsumer(&mocks.KafkaConsumerMock{})
}
