package qkafka

var consumerMock consumer = nil

func MockProducer(m producer) {
	producerInstance = m
}

func MockConsumer(m consumer) {
	consumerMock = m
}
