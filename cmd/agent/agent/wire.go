//go:build wireinject
// +build wireinject

package agent

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/agent/internal/deploy"
	"github.com/kappusuton-yon-tebaru/backend/cmd/agent/internal/deployenv"
	"github.com/kappusuton-yon-tebaru/backend/cmd/agent/internal/setting"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	sharedDeployEnv "github.com/kappusuton-yon-tebaru/backend/internal/deployenv"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/mongodb"
	"github.com/kappusuton-yon-tebaru/backend/internal/role"
	"github.com/kappusuton-yon-tebaru/backend/internal/user"
	"github.com/kappusuton-yon-tebaru/backend/internal/resource"
	"github.com/kappusuton-yon-tebaru/backend/internal/resourcerelationship"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type App struct {
	Logger           *logger.Logger
	Config           *config.Config
	SettingHandler   *setting.Handler
	DeployHandler    *deploy.Handler
	DeployEnvHandler *deployenv.Handler
}

func New(
	Logger *logger.Logger,
	Config *config.Config,
	SettingHandler *setting.Handler,
	DeployHandler *deploy.Handler,
	DeployEnvHandler *deployenv.Handler,
) *App {
	return &App{
		Logger,
		Config,
		SettingHandler,
		DeployHandler,
		DeployEnvHandler,
	}
}

func Initialize() (*App, error) {
	wire.Build(
		config.Load,
		logger.New,
		kubernetes.New,
		validator.New,
		setting.NewService,
		setting.NewHandler,
		deployenv.NewHandler,
		sharedDeployEnv.NewService,
		deploy.NewService,
		deploy.NewHandler,
		mongodb.NewMongoDB,
		resourcerelationship.NewRepository,
		resource.NewRepository,
		resource.NewService,
		role.NewRepository,
		user.NewRepository,
		New,
	)

	return &App{}, nil
}
