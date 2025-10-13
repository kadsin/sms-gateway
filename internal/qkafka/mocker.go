package qkafka

var producerMock Producer = nil

var consumerMock Consumer = nil

func MockProducer(m Producer) {
	producerMock = m
}

func MockConsumer(m Consumer) {
	consumerMock = m
}
