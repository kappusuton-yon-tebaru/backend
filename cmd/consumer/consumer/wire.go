//go:build wireinject
// +build wireinject

package builderconsumer

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/consumer/internal/build"
	"github.com/kappusuton-yon-tebaru/backend/cmd/consumer/internal/deploy"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/mongodb"
	"github.com/kappusuton-yon-tebaru/backend/internal/rmq"
)

type App struct {
	Logger        *logger.Logger
	Config        *config.Config
	KubeClient    *kubernetes.Kubernetes
	RmqClient     *rmq.BuilderRmq
	BuildHandler  *build.Handler
	DeployHandler *deploy.Handler
}

func New(
	Logger *logger.Logger,
	Config *config.Config,
	KubeClient *kubernetes.Kubernetes,
	RmqClient *rmq.BuilderRmq,
	BuildHandler *build.Handler,
	DeployHandler *deploy.Handler,
) *App {
	return &App{
		Logger,
		Config,
		KubeClient,
		RmqClient,
		BuildHandler,
		DeployHandler,
	}
}

func Initialize() (*App, error) {
	wire.Build(
		config.Load,
		logger.New,
		kubernetes.New,
		rmq.New,
		build.NewService,
		build.NewHandler,
		mongodb.NewMongoDB,
		job.NewRepository,
		job.NewService,
		deploy.NewHandler,
		deploy.NewService,
		New,
	)

	return &App{}, nil
}
