package utils

type Consumer interface {
	Update(data ...any)
}

type Producer interface {
	RegisterConsumer(consumer Consumer)
	notifyConsumers(data ...any)
	DeregisterConsumer(consumerToBeRemoved Consumer) error
}
