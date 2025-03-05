package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/backend"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/router"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"go.uber.org/zap"
)

//	@title			Snapping Service
//	@description	Snapping Service API Documentation
//	@version		1.0
//	@host			localhost:3000
//	@BasePath		/
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
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.RegisterRoutes(app)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				app.Logger.Fatal("panic occured while starting gin", zap.Any("recovered", r), zap.Stack("stacktrace"))
			}
		}()

		if err := r.Run(fmt.Sprintf("0.0.0.0:%v", app.Config.Backend.Port)); err != nil {
			app.Logger.Fatal("error occured while running gin", zap.Error(err))
		}
	}()

	<-utils.WaitForTermination(app.Logger, func() {
		if err := app.MongoDatabase.Client().Disconnect(context.Background()); err != nil {
			app.Logger.Error("error occured while disconnecting from mongodb", zap.Error(err))
		}
	})
}
