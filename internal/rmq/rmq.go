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

type BuilderRmq struct {
	Conn  *amqp091.Connection
	Ch    *amqp091.Channel
	Queue amqp091.Queue
}

func New(cfg *config.Config) (*BuilderRmq, error) {
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

	return &BuilderRmq{
		conn,
		ch,
		queue,
	}, nil
}

func (r *BuilderRmq) Publish(ctx context.Context, key string, body []byte) error {
	err := r.Ch.PublishWithContext(ctx, exchangeName, fmt.Sprintf("%s.%s", r.Queue.Name, key), false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})

	if err != nil {
		return err
	}

	return nil
}
