package container

import "github.com/kadsin/sms-gateway/internal/qkafka"

var mockKafkaConsumerInstance qkafka.Consumer = nil

func MockKafkaProducer(p qkafka.Producer) {
	kafkaProducerInstance = p
}

func MockKafkaConsumer(c qkafka.Consumer) {
	mockKafkaConsumerInstance = c
}
