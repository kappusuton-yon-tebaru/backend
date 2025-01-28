//go:build wireinject
// +build wireinject

package backend

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/greeting"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/image"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/resource"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/svcdeploy"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/user"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	sharedImage "github.com/kappusuton-yon-tebaru/backend/internal/image"
	"github.com/kappusuton-yon-tebaru/backend/internal/mongodb"
	sharedResource "github.com/kappusuton-yon-tebaru/backend/internal/resource"
	sharedSvcDeploy "github.com/kappusuton-yon-tebaru/backend/internal/svcdeploy"
	sharedUser "github.com/kappusuton-yon-tebaru/backend/internal/user"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	Config            *config.Config
	GreetingHandler   *greeting.Handler
	MongoClient       *mongo.Client
	ImageHandler      *image.Handler
	ServiceDeployment *svcdeploy.Handler
	UserHandler       *user.Handler
	ResourceHandler   *resource.Handler
}

func New(
	Config *config.Config,
	GreetingHandler *greeting.Handler,
	MongoClient *mongo.Client,
	ImageHandler *image.Handler,
	ServiceDeployment *svcdeploy.Handler,
	UserHandler *user.Handler,
	ResourceHandler *resource.Handler,
) *App {
	return &App{
		Config,
		GreetingHandler,
		MongoClient,
		ImageHandler,
		ServiceDeployment,
		UserHandler,
		ResourceHandler,
	}
}

func Initialize() (*App, error) {
	wire.Build(
		config.Load,
		greeting.New,
		mongodb.New,
		sharedImage.NewRepository,
		sharedImage.NewService,
		image.NewHandler,
		sharedSvcDeploy.NewRepository,
		sharedSvcDeploy.NewService,
		svcdeploy.NewHandler,
		sharedUser.NewRepository,
		sharedUser.NewService,
		user.NewHandler,
		sharedResource.NewRepository,
		sharedResource.NewService,
		resource.NewHandler,
		New,
	)

	return &App{}, nil
}
