//go:build wireinject
// +build wireinject

package podlogger

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/podlogger/internal/podevent"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/logging"
	"github.com/kappusuton-yon-tebaru/backend/internal/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	Logger        *logger.Logger
	Config        *config.Config
	KubeClient    *kubernetes.Kubernetes
	MongoDatabase *mongo.Database
	LogService    *logging.Service
	PodHandler    *podevent.Handler
}

func New(
	Logger *logger.Logger,
	Config *config.Config,
	KubeClient *kubernetes.Kubernetes,
	MongoDatabase *mongo.Database,
	LogService *logging.Service,
	PodHandler *podevent.Handler,
) *App {
	return &App{
		Logger,
		Config,
		KubeClient,
		MongoDatabase,
		LogService,
		PodHandler,
	}
}

func Initialize() (*App, error) {
	wire.Build(
		config.Load,
		logger.New,
		kubernetes.New,
		mongodb.NewMongoDB,
		logging.NewRepository,
		logging.NewService,
		podevent.NewHandler,
		New,
	)

	return &App{}, nil
}
