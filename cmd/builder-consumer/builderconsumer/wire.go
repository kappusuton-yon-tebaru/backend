//go:build wireinject
// +build wireinject

package builderconsumer

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/builder-consumer/internal/build"
	"github.com/kappusuton-yon-tebaru/backend/cmd/builder-consumer/internal/rmq"
	sharedBuild "github.com/kappusuton-yon-tebaru/backend/internal/build"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
)

type App struct {
	Logger       *logger.Logger
	Config       *config.Config
	KubeClient   *kubernetes.Kubernetes
	RmqClient    *rmq.Rmq
	BuildHandler *build.Handler
}

func New(
	Logger *logger.Logger,
	Config *config.Config,
	KubeClient *kubernetes.Kubernetes,
	RmqClient *rmq.Rmq,
	BuildHandler *build.Handler,
) *App {
	return &App{
		Logger,
		Config,
		KubeClient,
		RmqClient,
		BuildHandler,
	}
}

func Initialize() (*App, error) {
	wire.Build(
		config.Load,
		logger.New,
		kubernetes.New,
		rmq.New,
		sharedBuild.NewService,
		build.NewHandler,
		New,
	)

	return &App{}, nil
}
