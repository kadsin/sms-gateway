package qkafka

import "github.com/segmentio/kafka-go"

func NewProducer(c kafka.WriterConfig) Producer {
	if producerMock != nil {
		return producerMock
	}

	return &KafkaProducer{
		writer: kafka.NewWriter(c),
	}
}

func NewConsumer(c kafka.ReaderConfig) Consumer {
	if consumerMock != nil {
		return consumerMock
	}

	return &KafkaConsumer{
		reader: kafka.NewReader(c),
	}
}
