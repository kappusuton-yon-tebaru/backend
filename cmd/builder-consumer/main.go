package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kappusuton-yon-tebaru/backend/cmd/builder-consumer/builderconsumer"
)

func main() {
	app, err := builderconsumer.Initialize()
	if err != nil {
		panic(err)
	}

	fmt.Println("Builder consumer started")

	msgs, err := app.RmqClient.Ch.Consume(app.Config.BuilderConfig.QueueName, "agent-builder-consumer", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		err := app.BuildHandler.BuildImageHandler(msg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	<-waitForCleanup(app)
}

func waitForCleanup(app *builderconsumer.App) <-chan bool {
	done := make(chan bool, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig

		if err := app.RmqClient.Conn.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}

		if err := app.RmqClient.Ch.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}

		done <- true
	}()

	return done
}
