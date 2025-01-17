//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/agent/internal/greeting"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	shared_greeting "github.com/kappusuton-yon-tebaru/backend/internal/greeting"
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
		shared_greeting.New,
		New,
	)

	return &App{}, nil
}
