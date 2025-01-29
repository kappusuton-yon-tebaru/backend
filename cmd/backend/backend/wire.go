//go:build wireinject
// +build wireinject

package backend

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/greeting"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/image"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/projectrepository"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/resource"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/role"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/svcdeploy"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/svcdeployenv"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/user"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/usergroup"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	sharedImage "github.com/kappusuton-yon-tebaru/backend/internal/image"
	"github.com/kappusuton-yon-tebaru/backend/internal/mongodb"
	sharedProjectRepository "github.com/kappusuton-yon-tebaru/backend/internal/projectrepository"
	sharedResource "github.com/kappusuton-yon-tebaru/backend/internal/resource"
	sharedRole "github.com/kappusuton-yon-tebaru/backend/internal/role"
	sharedSvcDeploy "github.com/kappusuton-yon-tebaru/backend/internal/svcdeploy"
	sharedSvcDeployEnv "github.com/kappusuton-yon-tebaru/backend/internal/svcdeployenv"
	sharedUser "github.com/kappusuton-yon-tebaru/backend/internal/user"
	sharedUserGroup "github.com/kappusuton-yon-tebaru/backend/internal/usergroup"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	Config                   *config.Config
	GreetingHandler          *greeting.Handler
	MongoClient              *mongo.Client
	ImageHandler             *image.Handler
	ServiceDeployment        *svcdeploy.Handler
	ServiceDeploymentEnv     *svcdeployenv.Handler
	UserHandler              *user.Handler
	UserGroupHandler         *usergroup.Handler
	ResourceHandler          *resource.Handler
	RoleHandler              *role.Handler
	ProjectRepositoryHandler *projectrepository.Handler
}

func New(
	Config *config.Config,
	GreetingHandler *greeting.Handler,
	MongoClient *mongo.Client,
	ImageHandler *image.Handler,
	ServiceDeployment *svcdeploy.Handler,
	ServiceDeploymentEnv *svcdeployenv.Handler,
	UserHandler *user.Handler,
	UserGroupHandler *usergroup.Handler,
	ResourceHandler *resource.Handler,
	RoleHandler *role.Handler,
	ProjectRepositoryHandler *projectrepository.Handler,
) *App {
	return &App{
		Config,
		GreetingHandler,
		MongoClient,
		ImageHandler,
		ServiceDeployment,
		ServiceDeploymentEnv,
		UserHandler,
		UserGroupHandler,
		ResourceHandler,
		RoleHandler,
		ProjectRepositoryHandler,
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
		sharedSvcDeployEnv.NewRepository,
		sharedSvcDeployEnv.NewService,
		svcdeployenv.NewHandler,
		sharedUser.NewRepository,
		sharedUser.NewService,
		user.NewHandler,
		sharedUserGroup.NewRepository,
		sharedUserGroup.NewService,
		usergroup.NewHandler,
		sharedResource.NewRepository,
		sharedResource.NewService,
		resource.NewHandler,
		sharedRole.NewRepository,
		sharedRole.NewService,
		role.NewHandler,
		sharedProjectRepository.NewRepository,
		sharedProjectRepository.NewService,
		projectrepository.NewHandler,
		New,
	)

	return &App{}, nil
}
