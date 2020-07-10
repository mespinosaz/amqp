package queue

import (
	"github.com/streadway/amqp"
)

func NewFactory(ch *amqp.Channel) Factory {
	return &factory{
		ch: ch,
	}
}

type Factory interface {
	Create(n string) (Queue, error)
}

type factory struct {
	ch *amqp.Channel
}

func (f *factory) Create(n string) (Queue, error) {
	err := f.ch.Qos(10, 0, false)

	q, err := f.ch.QueueDeclare(
		n,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return NewQueue(f.ch, &q), nil
}
