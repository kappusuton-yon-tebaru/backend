//go:build wireinject
// +build wireinject

package agent

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/agent/internal/greeting"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	shared_greeting "github.com/kappusuton-yon-tebaru/backend/internal/greeting"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
)

type App struct {
	Logger          *logger.Logger
	Config          *config.Config
	GreetingHandler *greeting.Handler
}

func New(
	Logger *logger.Logger,
	Config *config.Config,
	GreetingHandler *greeting.Handler,
) *App {
	return &App{
		Logger,
		Config,
		GreetingHandler,
	}
}

func Initialize() (*App, error) {
	wire.Build(
		config.Load,
		logger.New,
		greeting.New,
		shared_greeting.New,
		New,
	)

	return &App{}, nil
}
