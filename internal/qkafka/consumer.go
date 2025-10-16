package qkafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func NewConsumer(c kafka.ReaderConfig) Consumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(c),
	}
}

type Consumer interface {
	Config() kafka.ReaderConfig
	FetchMessage(ctx context.Context) (kafka.Message, error)
	Commit(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type KafkaConsumer struct {
	reader *kafka.Reader
}

func (c *KafkaConsumer) Config() kafka.ReaderConfig {
	return c.reader.Config()
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
