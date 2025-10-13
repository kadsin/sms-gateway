package qkafka

import "github.com/segmentio/kafka-go"

var producerMock Producer = nil

func MockProducer(m Producer) {
	producerMock = m
}

func NewProducer(c kafka.WriterConfig) Producer {
	if producerMock != nil {
		return producerMock
	}

	return &KafkaProducer{
		writer: kafka.NewWriter(c),
	}
}
