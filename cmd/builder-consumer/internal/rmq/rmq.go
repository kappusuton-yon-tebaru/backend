package rmq

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

type Rmq struct {
	Conn  *amqp091.Connection
	Ch    *amqp091.Channel
	Queue amqp091.Queue
}

func New(cfg *config.Config) (*Rmq, error) {
	conn, err := amqp091.Dial(cfg.BuilderConfig.QueueUri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(cfg.BuilderConfig.QueueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &Rmq{
		conn,
		ch,
		queue,
	}, nil
}
