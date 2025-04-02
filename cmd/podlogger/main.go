package main

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/cmd/podlogger/podlogger"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"go.uber.org/zap"
)

func main() {
	app, err := podlogger.Initialize()
	if err != nil {
		panic(err)
	}

	app.Logger.Info("pod logger initalizing")

	podInformer, err := app.KubeClient.NewPodInformer(app.PodHandler, "watchlog=true")
	if err != nil {
		panic(err)
	}

	app.Logger.Info("pod logger ready to collect logs", zap.String("mode", string(app.Config.PodLogger.Mode)))

	<-utils.WaitForTermination(app.Logger, func() {
		app.Logger.Info("stopping pod informer")
		podInformer.Stop()

		app.Logger.Info("closing mongodb connection")
		if err := app.MongoDatabase.Client().Disconnect(context.Background()); err != nil {
			app.Logger.Error("error occured while disconnecting from mongodb", zap.Error(err))
		}
	})
}
