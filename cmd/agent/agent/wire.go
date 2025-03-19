//go:build wireinject
// +build wireinject

package agent

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/agent/internal/deploy"
	"github.com/kappusuton-yon-tebaru/backend/cmd/agent/internal/monitoring"
	"github.com/kappusuton-yon-tebaru/backend/cmd/agent/internal/setting"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	sharedDeployEnv "github.com/kappusuton-yon-tebaru/backend/internal/deployenv"
	"github.com/kappusuton-yon-tebaru/backend/internal/hub"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/mongodb"
	"github.com/kappusuton-yon-tebaru/backend/internal/resource"
	"github.com/kappusuton-yon-tebaru/backend/internal/resourcerelationship"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type App struct {
	Logger            *logger.Logger
	Config            *config.Config
	MonitoringHandler *monitoring.Handler
	SettingHandler    *setting.Handler
	DeployHandler     *deploy.Handler
}

func New(
	Logger *logger.Logger,
	Config *config.Config,
	MonitoringHandler *monitoring.Handler,
	SettingHandler *setting.Handler,
	DeployHandler *deploy.Handler,
) *App {
	return &App{
		Logger,
		Config,
		MonitoringHandler,
		SettingHandler,
		DeployHandler,
	}
}

func Initialize() (*App, error) {
	wire.Build(
		config.Load,
		logger.New,
		kubernetes.New,
		hub.New,
		monitoring.NewService,
		monitoring.NewHandler,
		validator.New,
		setting.NewService,
		setting.NewHandler,
		sharedDeployEnv.NewService,
		deploy.NewHandler,
		mongodb.NewMongoDB,
		resourcerelationship.NewRepository,
		resource.NewRepository,
		resource.NewService,
		New,
	)

	return &App{}, nil
}
