package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/cmd/agent/agent"
	"github.com/kappusuton-yon-tebaru/backend/cmd/agent/internal/router"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"go.uber.org/zap"
)

func main() {
	app, err := agent.Initialize()
	if err != nil {
		panic(err)
	}

	app.Logger.Info("agent initalizing")

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

	r.RegisterRoutes(app)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				app.Logger.Fatal("panic occured while starting gin", zap.Any("recovered", r), zap.Stack("stacktrace"))
			}
		}()

		app.Logger.Info("start serving request")

		if err := r.Run(fmt.Sprintf("0.0.0.0:%v", app.Config.Backend.Port)); err != nil {
			app.Logger.Fatal("error occured while running gin", zap.Error(err))
		}
	}()

	<-utils.WaitForTermination(app.Logger, func() {})
}

