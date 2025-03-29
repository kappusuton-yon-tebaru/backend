package rmq

import (
	"context"
	"fmt"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	exchangeName = "amq.topic"
)

type Rmq struct {
	logger    *logger.Logger
	conn      *amqp091.Connection
	ch        *amqp091.Channel
	queueUri  string
	queueName string
}

func New(cfg *config.Config, logger *logger.Logger) (*Rmq, error) {
	rmq := &Rmq{
		logger,
		nil,
		nil,
		cfg.ConsumerConfig.QueueUri,
		cfg.ConsumerConfig.OrganizationName,
	}

	if err := rmq.connect(); err != nil {
		return nil, err
	}

	go rmq.autoReconnect()

	return rmq, nil
}

func (r *Rmq) connect() error {
	conn, err := amqp091.Dial(r.queueUri)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	err = ch.Qos(1, 0, true)
	if err != nil {
		return err
	}

	queue, err := ch.QueueDeclare(r.queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.QueueBind(queue.Name, fmt.Sprintf("%s.*", queue.Name), exchangeName, false, nil)
	if err != nil {
		return err
	}

	r.conn = conn
	r.ch = ch

	r.logger.Info("rmq connection established", zap.Error(err))

	return nil
}

func (r *Rmq) autoReconnect() {
	for {
		sig := r.conn.NotifyClose(make(chan *amqp091.Error))

		for {
			err, ok := <-sig
			if !ok {
				break
			}

			r.logger.Error("rmq connection closed", zap.Error(err))

			if err := r.connect(); err != nil {
				r.logger.Error("cannot reestablish rmq connection", zap.Error(err))
				return
			}
		}
	}
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
	err := r.ch.PublishWithContext(ctx, exchangeName, fmt.Sprintf("%s.%s", r.queueName, key), false, false, amqp091.Publishing{
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
