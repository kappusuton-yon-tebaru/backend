package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/backend"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/router"
)

func main() {
	app, err := backend.Initialize()
	if err != nil {
		panic(err)
	}

	defer app.MongoClient.Disconnect(context.Background())

	r := router.New()
	config := cors.Config{
		AllowOrigins:     []string{"http://example.com", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}

	r.Use(cors.New(config))
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
