//go:build wireinject
// +build wireinject

package backend

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/greeting"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/image"
	"github.com/kappusuton-yon-tebaru/backend/internal/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	Config          *config.Config
	GreetingHandler *greeting.Handler
	MongoClient     *mongo.Client
	ImageRepo       *image.Repository
}

func New(
	Config *config.Config,
	GreetingHandler *greeting.Handler,
	MongoClient *mongo.Client,
	ImageRepo *image.Repository,
) *App {
	return &App{
		Config,
		GreetingHandler,
		MongoClient,
		ImageRepo,
	}
}

func Initialize() (*App, error) {
	wire.Build(
		config.Load,
		greeting.New,
		mongodb.New,
		image.NewRepository,
		New,
	)

	return &App{}, nil
}
