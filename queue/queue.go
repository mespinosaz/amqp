package queue

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func NewQueue(ch *amqp.Channel, q *amqp.Queue) Queue {
	return &queue{
		ch: ch,
		q:  q,
	}
}

type Queue interface {
	Publish(message interface{}) error
	Consume(tag string, f func(amqp.Delivery) error)
}

type queue struct {
	ch *amqp.Channel
	q  *amqp.Queue
}

func (iq *queue) Publish(m interface{}) error {
	b, err := json.Marshal(m)

	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return err
	}

	err = iq.ch.Publish(
		"",
		iq.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        b,
		})

	if err != nil {
		return err
	}

	log.Printf("Published in queue %s message %v", iq.q.Name, m)

	return nil
}

func (iq *queue) Consume(tag string, f func(amqp.Delivery) error) {
	log.Printf("Consuming messages from queue %s", iq.q.Name)
	msgs, err := iq.ch.Consume(
		iq.q.Name,
		tag,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		err = f(d)

		if err != nil {
			log.Printf("Error: %s", err)
			d.Nack(true, true)
			return
		}

		d.Ack(true)
	}
}
