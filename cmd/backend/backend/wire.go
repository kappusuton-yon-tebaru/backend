//go:build wireinject
// +build wireinject

package backend

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/greeting"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/image"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	shared_image "github.com/kappusuton-yon-tebaru/backend/internal/image"
	"github.com/kappusuton-yon-tebaru/backend/internal/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	Config          *config.Config
	GreetingHandler *greeting.Handler
	MongoClient     *mongo.Client
	ImageHandler    *image.Handler
}

func New(
	Config *config.Config,
	GreetingHandler *greeting.Handler,
	MongoClient *mongo.Client,
	ImageHandler *image.Handler,
) *App {
	return &App{
		Config,
		GreetingHandler,
		MongoClient,
		ImageHandler,
	}
}

func Initialize() (*App, error) {
	wire.Build(
		config.Load,
		greeting.New,
		mongodb.New,
		shared_image.NewRepository,
		shared_image.NewService,
		image.NewHandler,
		New,
	)

	return &App{}, nil
}
