package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kappusuton-yon-tebaru/backend/cmd/builder-consumer/builderconsumer"
	"go.uber.org/zap"
)

func main() {
	app, err := builderconsumer.Initialize()
	if err != nil {
		panic(err)
	}

	app.Logger.Info("builder consumer initalizing")

	app.Logger.Info("connecting to rmq", zap.String("queue_name", app.Config.BuilderConfig.QueueUri))
	msgs, err := app.RmqClient.Ch.Consume(app.Config.BuilderConfig.QueueName, "agent-builder-consumer", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	app.Logger.Info("builder consumer is ready to consume message")
	for msg := range msgs {
		func() {
			defer func() {
				if r := recover(); r != nil {
					app.Logger.Error("panic occured", zap.Any("recovered", r), zap.Stack("stacktrace"))
				}
			}()

			err := app.BuildHandler.BuildImageHandler(msg)
			if err != nil {
				app.Logger.Error("error occured while handling message", zap.Error(err))
			}
		}()
	}

	<-waitForCleanup(app)
}

func waitForCleanup(app *builderconsumer.App) <-chan bool {
	done := make(chan bool, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		app.Logger.Info("cleaning up before shutdown")

		app.Logger.Info("closing rmq channel")
		if err := app.RmqClient.Ch.Close(); err != nil {
			app.Logger.Error("error occured while closing rmq channel", zap.Error(err))
		}

		app.Logger.Info("closing rmq connection")
		if err := app.RmqClient.Conn.Close(); err != nil {
			app.Logger.Error("error occured while closing rmq connection", zap.Error(err))
		}

		app.Logger.Info("shutting down")

		done <- true
	}()

	return done
}
