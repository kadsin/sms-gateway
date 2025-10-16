package container

import (
	"github.com/kadsin/sms-gateway/internal/qkafka"
	"github.com/kadsin/sms-gateway/internal/sms"
)

func MockKafkaProducer(p qkafka.Producer) {
	kafkaProducerInstance = p
}

var mockKafkaConsumerInstance qkafka.Consumer = nil

func MockKafkaConsumer(c qkafka.Consumer) {
	mockKafkaConsumerInstance = c
}

var mockSmsProviderInstance sms.Provider = nil

func MockSmsProvider(p sms.Provider) {
	mockSmsProviderInstance = p
}
