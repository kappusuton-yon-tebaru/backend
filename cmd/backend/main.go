package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/backend"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/router"
)

func main() {
	app, err := backend.Initialize()
	if err != nil {
		panic(err)
	}

	r := router.New()
	r.RegisterRoutes(app)

	go func() {
		if err := r.Run(fmt.Sprintf("0.0.0.0:%v", app.Config.Backend.Port)); err != nil {
			panic(err)
		}
	}()

	<-waitForCleanup(app)
}

func waitForCleanup(app *backend.App) <-chan bool {
	done := make(chan bool, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig

		if err := app.MongoClient.Disconnect(context.Background()); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}

		done <- true
	}()

	return done
}
