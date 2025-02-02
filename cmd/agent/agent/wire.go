//go:build wireinject
// +build wireinject

package agent

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/agent/internal/monitoring"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/hub"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	shared_monitoring "github.com/kappusuton-yon-tebaru/backend/internal/monitoring"
)

type App struct {
	Logger            *logger.Logger
	Config            *config.Config
	MonitoringHandler *monitoring.Handler
}

func New(
	Logger *logger.Logger,
	Config *config.Config,
	MonitoringHandler *monitoring.Handler,
) *App {
	return &App{
		Logger,
		Config,
		MonitoringHandler,
	}
}

func Initialize() (*App, error) {
	wire.Build(
		config.Load,
		logger.New,
		kubernetes.New,
		hub.New,
		shared_monitoring.NewService,
		monitoring.NewHandler,
		New,
	)

	return &App{}, nil
}
