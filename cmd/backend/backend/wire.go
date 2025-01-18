//go:build wireinject
// +build wireinject

package backend

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/greeting"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
)

type App struct {
	Config          *config.Config
	GreetingHandler *greeting.Handler
}

func New(
	Config *config.Config,
	GreetingHandler *greeting.Handler,
) *App {
	return &App{
		Config,
		GreetingHandler,
	}
}

func Initialize() (*App, error) {
	wire.Build(
		config.Load,
		greeting.New,
		New,
	)

	return &App{}, nil
}
