package rmq

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

type BuilderRmq struct {
	Conn      *amqp091.Connection
	Ch        *amqp091.Channel
	Queue     amqp091.Queue
	QueueName string
}

func New(cfg *config.Config) (*BuilderRmq, error) {
	conn, err := amqp091.Dial(cfg.BuilderConfig.QueueUri)
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

	queue, err := ch.QueueDeclare(cfg.BuilderConfig.QueueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &BuilderRmq{
		conn,
		ch,
		queue,
		cfg.BuilderConfig.QueueName,
	}, nil
}

func (r *BuilderRmq) Publish(ctx context.Context, body []byte) error {
	err := r.Ch.PublishWithContext(ctx, "", r.QueueName, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})

	if err != nil {
		return err
	}

	return nil
}
