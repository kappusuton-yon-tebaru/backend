//go:build wireinject
// +build wireinject

package backend

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/build"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/greeting"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/image"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/monitoring"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/permission"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/projectenv"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/projectrepository"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/regproviders"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/resource"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/resourcerelationship"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/reverseproxy"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/role"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/rolepermission"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/roleusergroup"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/svcdeploy"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/svcdeployenv"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/user"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/usergroup"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	sharedImage "github.com/kappusuton-yon-tebaru/backend/internal/image"
	sharedJob "github.com/kappusuton-yon-tebaru/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/mongodb"
	sharedPermission "github.com/kappusuton-yon-tebaru/backend/internal/permission"
	sharedProjectEnvironment "github.com/kappusuton-yon-tebaru/backend/internal/projectenv"
	sharedProjectRepository "github.com/kappusuton-yon-tebaru/backend/internal/projectrepository"
	sharedRegProviders "github.com/kappusuton-yon-tebaru/backend/internal/regproviders"
	sharedResource "github.com/kappusuton-yon-tebaru/backend/internal/resource"
	sharedResourceRelationship "github.com/kappusuton-yon-tebaru/backend/internal/resourcerelationship"
	"github.com/kappusuton-yon-tebaru/backend/internal/rmq"
	sharedRole "github.com/kappusuton-yon-tebaru/backend/internal/role"
	sharedRolePermission "github.com/kappusuton-yon-tebaru/backend/internal/rolepermission"
	sharedRoleUserGroup "github.com/kappusuton-yon-tebaru/backend/internal/roleusergroup"
	sharedSvcDeploy "github.com/kappusuton-yon-tebaru/backend/internal/svcdeploy"
	sharedSvcDeployEnv "github.com/kappusuton-yon-tebaru/backend/internal/svcdeployenv"
	sharedUser "github.com/kappusuton-yon-tebaru/backend/internal/user"
	sharedUserGroup "github.com/kappusuton-yon-tebaru/backend/internal/usergroup"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	Logger                      *logger.Logger
	Config                      *config.Config
	GreetingHandler             *greeting.Handler
	MongoDatabase               *mongo.Database
	ImageHandler                *image.Handler
	ServiceDeployment           *svcdeploy.Handler
	ServiceDeploymentEnv        *svcdeployenv.Handler
	UserHandler                 *user.Handler
	UserGroupHandler            *usergroup.Handler
	ResourceHandler             *resource.Handler
	RoleHandler                 *role.Handler
	PermissionHandler           *permission.Handler
	RolePermissionHandler       *rolepermission.Handler
	RoleUserGroupHandler        *roleusergroup.Handler
	ProjectRepositoryHandler    *projectrepository.Handler
	ResourceRelationshipHandler *resourcerelationship.Handler
	JobHandler                  *job.Handler
	RegisterProviderHandler     *regproviders.Handler
	ProjectEnvironmentHandler   *projectenv.Handler
	BuildHandler                *build.Handler
	MonitoringHandler           *monitoring.Handler
	ReverseProxyHandler         *reverseproxy.ReverseProxy
}

func New(
	Logger *logger.Logger,
	Config *config.Config,
	GreetingHandler *greeting.Handler,
	MongoDatabase *mongo.Database,
	ImageHandler *image.Handler,
	ServiceDeployment *svcdeploy.Handler,
	ServiceDeploymentEnv *svcdeployenv.Handler,
	UserHandler *user.Handler,
	UserGroupHandler *usergroup.Handler,
	ResourceHandler *resource.Handler,
	RoleHandler *role.Handler,
	PermissionHandler *permission.Handler,
	RolePermissionHandler *rolepermission.Handler,
	RoleUserGroupHandler *roleusergroup.Handler,
	ProjectRepositoryHandler *projectrepository.Handler,
	ResourceRelationshipHandler *resourcerelationship.Handler,
	JobHandler *job.Handler,
	RegisterProviderHandler *regproviders.Handler,
	ProjectEnvironmentHandler *projectenv.Handler,
	BuildHandler *build.Handler,
	MonitoringHandler *monitoring.Handler,
	ReverseProxyHandler *reverseproxy.ReverseProxy,
) *App {
	return &App{
		Logger,
		Config,
		GreetingHandler,
		MongoDatabase,
		ImageHandler,
		ServiceDeployment,
		ServiceDeploymentEnv,
		UserHandler,
		UserGroupHandler,
		ResourceHandler,
		RoleHandler,
		PermissionHandler,
		RolePermissionHandler,
		RoleUserGroupHandler,
		ProjectRepositoryHandler,
		ResourceRelationshipHandler,
		JobHandler,
		RegisterProviderHandler,
		ProjectEnvironmentHandler,
		BuildHandler,
		MonitoringHandler,
		ReverseProxyHandler,
	}
}

func Initialize() (*App, error) {
	wire.Build(
		config.Load,
		logger.New,
		greeting.New,
		mongodb.NewMongoDB,
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
		sharedRoleUserGroup.NewRepository,
		sharedRoleUserGroup.NewService,
		roleusergroup.NewHandler,
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
		validator.New,
		rmq.New,
		build.NewService,
		build.NewHandler,
		monitoring.NewHandler,
		reverseproxy.New,
		New,
	)

	return &App{}, nil
}
