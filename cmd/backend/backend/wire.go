//go:build wireinject
// +build wireinject

package backend

import (
	"github.com/google/wire"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/auth"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/build"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/deploy"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/githubapi"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/greeting"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/image"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/monitoring"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/projectenv"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/projectrepository"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/regproviders"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/resource"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/resourcerelationship"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/reverseproxy"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/role"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/svcdeploy"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/svcdeployenv"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/user"
	// sharedECR "github.com/kappusuton-yon-tebaru/backend/internal/ecr"
	// sharedDockerHub "github.com/kappusuton-yon-tebaru/backend/internal/dockerhub"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/dockerhub"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/ecr"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/middleware"

	sharedGithubAPI "github.com/kappusuton-yon-tebaru/backend/internal/githubapi"

	sharedAuth "github.com/kappusuton-yon-tebaru/backend/internal/auth"
	sharedImage "github.com/kappusuton-yon-tebaru/backend/internal/image"
	sharedJob "github.com/kappusuton-yon-tebaru/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/mongodb"
	sharedProjectEnvironment "github.com/kappusuton-yon-tebaru/backend/internal/projectenv"
	sharedProjectRepository "github.com/kappusuton-yon-tebaru/backend/internal/projectrepository"
	sharedRegProviders "github.com/kappusuton-yon-tebaru/backend/internal/regproviders"
	sharedResource "github.com/kappusuton-yon-tebaru/backend/internal/resource"
	sharedResourceRelationship "github.com/kappusuton-yon-tebaru/backend/internal/resourcerelationship"
	"github.com/kappusuton-yon-tebaru/backend/internal/rmq"
	sharedRole "github.com/kappusuton-yon-tebaru/backend/internal/role"
	sharedSvcDeploy "github.com/kappusuton-yon-tebaru/backend/internal/svcdeploy"
	sharedSvcDeployEnv "github.com/kappusuton-yon-tebaru/backend/internal/svcdeployenv"
	sharedUser "github.com/kappusuton-yon-tebaru/backend/internal/user"
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
	ResourceHandler             *resource.Handler
	RoleHandler                 *role.Handler
	ProjectRepositoryHandler    *projectrepository.Handler
	ResourceRelationshipHandler *resourcerelationship.Handler
	JobHandler                  *job.Handler
	RegisterProviderHandler     *regproviders.Handler
	ProjectEnvironmentHandler   *projectenv.Handler
	ECRHandler                  *ecr.Handler
	DockerHubHandler            *dockerhub.Handler
	BuildHandler                *build.Handler
	DeployHandler               *deploy.Handler
	MonitoringHandler           *monitoring.Handler
	ReverseProxyHandler         *reverseproxy.ReverseProxy
	GithubAPIHandler            *githubapi.Handler
	AuthHandler                 *auth.Handler
	Middleware                  *middleware.Middleware
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
	ResourceHandler *resource.Handler,
	RoleHandler *role.Handler,
	ProjectRepositoryHandler *projectrepository.Handler,
	ResourceRelationshipHandler *resourcerelationship.Handler,
	JobHandler *job.Handler,
	RegisterProviderHandler *regproviders.Handler,
	ProjectEnvironmentHandler *projectenv.Handler,
	ECRHandler *ecr.Handler,
	DockerHubHandler *dockerhub.Handler,
	BuildHandler *build.Handler,
	MonitoringHandler *monitoring.Handler,
	ReverseProxyHandler *reverseproxy.ReverseProxy,
	GithubAPIHandler *githubapi.Handler,
	DeployHandler *deploy.Handler,
	AuthHandler *auth.Handler,
	Middleware *middleware.Middleware,
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
		ResourceHandler,
		RoleHandler,
		ProjectRepositoryHandler,
		ResourceRelationshipHandler,
		JobHandler,
		RegisterProviderHandler,
		ProjectEnvironmentHandler,
		ECRHandler,
		DockerHubHandler,
		BuildHandler,
		DeployHandler,
		MonitoringHandler,
		ReverseProxyHandler,
		GithubAPIHandler,
		AuthHandler,
		Middleware,
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
		sharedResource.NewRepository,
		sharedResource.NewService,
		resource.NewHandler,
		sharedRole.NewRepository,
		sharedRole.NewService,
		role.NewHandler,
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
		ecr.NewECRRepository,
		ecr.NewService,
		ecr.NewHandler,
		dockerhub.NewDockerHubRepository,
		dockerhub.NewService,
		dockerhub.NewHandler,
		validator.New,
		rmq.New,
		build.NewService,
		build.NewHandler,
		monitoring.NewHandler,
		reverseproxy.New,
		sharedGithubAPI.NewRepository,
		sharedGithubAPI.NewService,
		githubapi.NewHandler,
		deploy.NewHandler,
		deploy.NewService,
		auth.NewHandler,
		sharedAuth.NewService,
		sharedAuth.NewRepository,
		middleware.NewMiddleware,
		New,
	)

	return &App{}, nil
}
