package mocks

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumerMock struct {
	Topic   string
	counter int64
}

func (c *KafkaConsumerMock) FetchMessage(ctx context.Context) (kafka.Message, error) {
	kafkaTopicsMockMux.Lock()
	defer kafkaTopicsMockMux.Unlock()

	t, ok := kafkaTopicsMock[c.Topic]
	if !ok {
		return kafka.Message{}, fmt.Errorf("the topic %+v isn't exists", c.Topic)
	}

	if len(t) < 1 {
		return kafka.Message{}, errors.New("the topic is empty")
	}

	return t[c.counter], nil
}

func (c *KafkaConsumerMock) Commit(ctx context.Context, msgs ...kafka.Message) error {
	slices.SortFunc(msgs, func(a, b kafka.Message) int { return int(b.Offset - a.Offset) })
	lastMessage := msgs[len(msgs)-1]

	c.counter = lastMessage.Offset + 1

	return nil
}

func (c *KafkaConsumerMock) Close() error {
	c.counter = 0

	return nil
}
