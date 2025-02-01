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
)

type App struct {
	Config       *config.Config
	KubeClient   *kubernetes.Kubernetes
	RmqClient    *rmq.Rmq
	BuildHandler *build.Handler
}

func New(
	Config *config.Config,
	KubeClient *kubernetes.Kubernetes,
	RmqClient *rmq.Rmq,
	BuildHandler *build.Handler,
) *App {
	return &App{
		Config,
		KubeClient,
		RmqClient,
		BuildHandler,
	}
}

func Initialize() (*App, error) {
	wire.Build(
		config.Load,
		kubernetes.New,
		rmq.New,
		sharedBuild.NewService,
		build.NewHandler,
		New,
	)

	return &App{}, nil
}
