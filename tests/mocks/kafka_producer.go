package mocks

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaProducerMock struct {
	counter int64
}

func (p *KafkaProducerMock) SendMessage(ctx context.Context, m kafka.Message) error {
	kafkaTopicsMockMux.Lock()
	defer kafkaTopicsMockMux.Unlock()

	m.Offset = p.counter
	p.counter = p.counter + 1

	if _, ok := kafkaTopicsMock[m.Topic]; !ok {
		kafkaTopicsMock[m.Topic] = []kafka.Message{}
	}

	kafkaTopicsMock[m.Topic] = append(kafkaTopicsMock[m.Topic], m)
	return nil
}

func (p *KafkaProducerMock) Close() error {
	kafkaTopicsMock = map[string][]kafka.Message{}
	p.counter = 0

	return nil
}
