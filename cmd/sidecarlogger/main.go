package main

import (
	"bufio"
	"context"
	"fmt"

	"github.com/kappusuton-yon-tebaru/backend/cmd/sidecarlogger/sidecarlogger"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"go.uber.org/zap"
)

func main() {
	app, err := sidecarlogger.Initialize()
	if err != nil {
		panic(err)
	}

	go func() {
		podClient := app.KubeClient.NewPodClient(app.Config.SidecarLogger.Namespace)

		serviceName := app.Config.SidecarLogger.ServiceName

		reader, err := podClient.GetLog(serviceName, kubernetes.WithContainer(serviceName)).Stream(context.Background())
		if err != nil {
			panic(err)
		}
		defer reader.Close()

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()

			fmt.Println(line)
		}
	}()

	<-utils.WaitForTermination(app.Logger, func() {
		app.Logger.Info("closing rmq connection")
		if err := app.MongoDatabase.Client().Disconnect(context.Background()); err != nil {
			app.Logger.Error("error occured while disconnecting from mongodb", zap.Error(err))
		}
	})
}
