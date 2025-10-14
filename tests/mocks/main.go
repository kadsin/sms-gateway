package mocks

import (
	"sync"

	"github.com/segmentio/kafka-go"
)

var kafkaTopicsMock map[string][]kafka.Message = map[string][]kafka.Message{}
var kafkaTopicsMockMux sync.Mutex
