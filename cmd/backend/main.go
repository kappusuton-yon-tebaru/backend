package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/backend"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/router"
	"go.uber.org/zap"
)

func main() {
	app, err := backend.Initialize()
	if err != nil {
		panic(err)
	}

	r := router.New(app.Config)

	r.Use(func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		app.Logger.Info("request",
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.Int("status", ctx.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", ctx.ClientIP()),
		)
	})

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://example.com", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

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
