package rmq

import (
	"context"
	"fmt"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName = "amq.topic"
)

type Rmq struct {
	conn  *amqp091.Connection
	ch    *amqp091.Channel
	queue amqp091.Queue
}

func New(cfg *config.Config) (*Rmq, error) {
	conn, err := amqp091.Dial(cfg.ConsumerConfig.QueueUri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.Qos(1, 0, true)
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(cfg.ConsumerConfig.OrganizationName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(queue.Name, fmt.Sprintf("%s.*", queue.Name), exchangeName, false, nil)
	if err != nil {
		panic(err)
	}

	return &Rmq{
		conn,
		ch,
		queue,
	}, nil
}

func (r *Rmq) Close() error {
	if err := r.ch.Close(); err != nil {
		return err
	}

	if err := r.conn.Close(); err != nil {
		return err
	}

	return nil
}

func (r *Rmq) Publish(ctx context.Context, key string, body []byte) error {
	err := r.ch.PublishWithContext(ctx, exchangeName, fmt.Sprintf("%s.%s", r.queue.Name, key), false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *Rmq) Consume(queueName, consumerName string) (<-chan amqp091.Delivery, error) {
	msgs, err := r.ch.Consume(queueName, consumerName, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
