package qkafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func NewProducer(c kafka.WriterConfig) Producer {
	return &KafkaProducer{
		writer: kafka.NewWriter(c),
	}
}

type Producer interface {
	SendMessage(ctx context.Context, m kafka.Message) error
	Close() error
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func (p *KafkaProducer) SendMessage(ctx context.Context, m kafka.Message) error {
	return p.writer.WriteMessages(ctx, m)
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
