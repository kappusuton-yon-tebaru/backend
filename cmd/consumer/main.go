package main

import (
	"fmt"
	"strings"

	"github.com/kappusuton-yon-tebaru/backend/cmd/consumer/consumer"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"go.uber.org/zap"
)

func main() {
	app, err := builderconsumer.Initialize()
	if err != nil {
		panic(err)
	}

	app.Logger.Info("builder consumer initalizing")

	app.Logger.Info("connecting to rmq", zap.String("queue_name", app.Config.ConsumerConfig.OrganizationName), zap.String("queue_uri", app.Config.ConsumerConfig.QueueUri))
	msgs, err := app.RmqClient.Ch.Consume(app.Config.ConsumerConfig.OrganizationName, "agent-consumer", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		app.Logger.Info("builder consumer is ready to consume message")
		for msg := range msgs {
			func() {
				defer func() {
					if r := recover(); r != nil {
						app.Logger.Error("panic occured", zap.Any("recovered", r), zap.Stack("stacktrace"))
					}
				}()

				key := strings.TrimPrefix(msg.RoutingKey, fmt.Sprintf("%s.", app.Config.ConsumerConfig.OrganizationName))

				switch key {
				case enum.BuildContextRoutingKey:
					app.BuildHandler.BuildHandler(msg)
				case enum.DeployContextRoutingKey:
					app.DeployHandler.DeployHandler(msg)
				}
			}()

			if err := msg.Ack(false); err != nil {
				app.Logger.Error("error occured while acking", zap.Error(err))
			}
		}
	}()

	<-utils.WaitForTermination(app.Logger, func() {
		app.Logger.Info("closing rmq channel")
		if err := app.RmqClient.Ch.Close(); err != nil {
			app.Logger.Error("error occured while closing rmq channel", zap.Error(err))
		}

		app.Logger.Info("closing rmq connection")
		if err := app.RmqClient.Conn.Close(); err != nil {
			app.Logger.Error("error occured while closing rmq connection", zap.Error(err))
		}
	})
}
