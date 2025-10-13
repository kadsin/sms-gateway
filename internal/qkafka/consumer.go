package qkafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func NewConsumer(c kafka.ReaderConfig) consumer {
	if consumerMock != nil {
		return consumerMock
	}

	return &KafkaConsumer{
		reader: kafka.NewReader(c),
	}
}

type consumer interface {
	FetchMessage(ctx context.Context) (kafka.Message, error)
	Commit(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type KafkaConsumer struct {
	reader *kafka.Reader
}

func (c *KafkaConsumer) FetchMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.FetchMessage(ctx)
}

func (c *KafkaConsumer) Commit(ctx context.Context, msgs ...kafka.Message) error {
	return c.reader.CommitMessages(ctx, msgs...)
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
