//go:build wireinject
// +build wireinject

package backend

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/greeting"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/image"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/projectrepository"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/resource"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/resourcerelationship"
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
	sharedResourceRelationship "github.com/kappusuton-yon-tebaru/backend/internal/resourcerelationship"
	sharedRole "github.com/kappusuton-yon-tebaru/backend/internal/role"
	sharedSvcDeploy "github.com/kappusuton-yon-tebaru/backend/internal/svcdeploy"
	sharedSvcDeployEnv "github.com/kappusuton-yon-tebaru/backend/internal/svcdeployenv"
	sharedUser "github.com/kappusuton-yon-tebaru/backend/internal/user"
	sharedUserGroup "github.com/kappusuton-yon-tebaru/backend/internal/usergroup"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/permission"
	sharedPermission "github.com/kappusuton-yon-tebaru/backend/internal/permission"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/regproviders"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/projectenv"
	sharedJob "github.com/kappusuton-yon-tebaru/backend/internal/job"
	sharedRegProviders "github.com/kappusuton-yon-tebaru/backend/internal/regproviders"
	sharedProjectEnvironment "github.com/kappusuton-yon-tebaru/backend/internal/projectenv"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/rolepermission"
	sharedRolePermission "github.com/kappusuton-yon-tebaru/backend/internal/rolepermission"
)

type App struct {
	Config                      *config.Config
	GreetingHandler             *greeting.Handler
	MongoClient                 *mongo.Client
	ImageHandler                *image.Handler
	ServiceDeployment           *svcdeploy.Handler
	ServiceDeploymentEnv        *svcdeployenv.Handler
	UserHandler                 *user.Handler
	UserGroupHandler            *usergroup.Handler
	ResourceHandler             *resource.Handler
	RoleHandler                 *role.Handler
	PermissionHandler           *permission.Handler
	RolePermissionHandler       *rolepermission.Handler
	ProjectRepositoryHandler    *projectrepository.Handler
	ResourceRelationshipHandler *resourcerelationship.Handler
	JobHandler				  	*job.Handler
	RegisterProviderHandler		*regproviders.Handler
	ProjectEnvironmentHandler	*projectenv.Handler
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
	PermissionHandler *permission.Handler,
	RolePermissionHandler *rolepermission.Handler,
	ProjectRepositoryHandler *projectrepository.Handler,
	ResourceRelationshipHandler *resourcerelationship.Handler,
	JobHandler *job.Handler,
	RegisterProviderHandler *regproviders.Handler,
	ProjectEnvironmentHandler *projectenv.Handler,
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
		PermissionHandler,
		RolePermissionHandler,
		ProjectRepositoryHandler,
		ResourceRelationshipHandler,
		JobHandler,
		RegisterProviderHandler,
		ProjectEnvironmentHandler,
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
		sharedPermission.NewRepository,
		sharedPermission.NewService,
		permission.NewHandler,
		sharedRolePermission.NewRepository,
		sharedRolePermission.NewService,
		rolepermission.NewHandler,
		sharedProjectRepository.NewRepository,
		sharedProjectRepository.NewService,
		projectrepository.NewHandler,
		sharedResourceRelationship.NewRepository,
		sharedResourceRelationship.NewService,
		resourcerelationship.NewHandler,
		sharedJob.NewRepository,
		sharedJob.NewService,
		job.NewHandler,
		sharedRegProviders.NewRepository,
		sharedRegProviders.NewService,
		regproviders.NewHandler,
		sharedProjectEnvironment.NewRepository,
		sharedProjectEnvironment.NewService,
		projectenv.NewHandler,
		New,
	)

	return &App{}, nil
}
