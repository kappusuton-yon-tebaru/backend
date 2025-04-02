package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/backend"
	_ "github.com/kappusuton-yon-tebaru/backend/cmd/backend/docs"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	*gin.Engine
}

func New(cfg *config.Config) *Router {
	if cfg.Development {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	return &Router{
		r,
	}
}

func (r *Router) RegisterRoutes(app *backend.App) {
	if app.Config.Development {
		r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// example for authenticated route
	// authenticated := r.Group("/", app.Middleware.Authentication())
	// authenticated.GET("/", app.GreetingHandler.Greeting)

	r.GET("/users", app.UserHandler.GetAllUsers)
	r.DELETE("/users/:id", app.UserHandler.DeleteUserById)
	r.POST("/users/:user_id/roles/:role_id", app.UserHandler.AddRole)
	r.PUT("/users/:user_id/roles/:role_id", app.UserHandler.RemoveRole)
	authenticated.GET("/users/me/permissions", app.RoleHandler.GetUserPermissions)

	r.POST("/auth/register", app.AuthHandler.Register)
	r.POST("/auth/login", app.AuthHandler.Login)
	r.POST("/auth/logout", app.AuthHandler.Logout)

	r.GET("/resources", app.ResourceHandler.GetAllResources)
	r.GET("/resources/:id", app.ResourceHandler.GetResourceByID)
	r.GET("/resources/children/:parent_id", app.ResourceHandler.GetChildrenResourcesByParentID)
	authenticated.POST("/resources", app.ResourceHandler.CreateResource) // ?parent_id={id}
	r.PUT("/resources/:id", app.ResourceHandler.UpdateResource)
	r.DELETE("/resources/:id", app.ResourceHandler.DeleteResource)
	r.DELETE("/resources/cascade/:id", app.ResourceHandler.CascadeDeleteResource)

	r.GET("/roles", app.RoleHandler.GetAllRoles)
	r.POST("/roles", app.RoleHandler.CreateRole)
	r.PUT("/roles/:role_id", app.RoleHandler.UpdateRole)
	r.DELETE("/roles/:role_id", app.RoleHandler.DeleteRoleById)
	r.POST("/roles/:role_id/permissions", app.RoleHandler.AddPermission)
	r.PUT("/roles/:role_id/permissions/:perm_id", app.RoleHandler.UpdatePermission)
	r.DELETE("/roles/:role_id/permissions/:perm_id", app.RoleHandler.DeletePermission)

	r.GET("/projrepos", app.ProjectRepositoryHandler.GetAllProjectRepositories)
	r.GET("/projrepos/project/:project_id", app.ProjectRepositoryHandler.GetProjectRepositoryByProjectId)
	r.POST("/projrepos/:id", app.ProjectRepositoryHandler.CreateProjectRepository)
	r.PATCH("/projrepos/:id", app.ProjectRepositoryHandler.UpdateProjectRepositoryRegistryProvider)
	r.DELETE("/projrepos/:id", app.ProjectRepositoryHandler.DeleteProjectRepository)

	r.GET("/resourcerelas", app.ResourceRelationshipHandler.GetAllResourceRelationships)
	r.GET("/resourcerelas/:parent_id", app.ResourceRelationshipHandler.GetChildrenResourceRelationshipByParentID)
	r.POST("/resourcerelas", app.ResourceRelationshipHandler.CreateResourceRelationship)
	r.DELETE("/resourcerelas/:id", app.ResourceRelationshipHandler.DeleteResourceRelationship)

	r.GET("/jobs", app.JobHandler.GetAllJobParents)
	r.GET("/jobs/:id/parent", app.JobHandler.GetAllJobsByParentId)
	r.GET("/jobs/:id", app.JobHandler.GetJobById)
	r.GET("/jobs/:id/log", app.JobHandler.GetJobLog)

	r.GET("/regproviders", app.RegisterProviderHandler.GetAllRegProviders)
	r.GET("/regproviders/unbind", app.RegisterProviderHandler.GetAllRegProvidersWithoutProject)
	r.GET("/regproviders/:id", app.RegisterProviderHandler.GetRegProviderById)
	r.POST("/regproviders", app.RegisterProviderHandler.CreateRegProvider)
	r.DELETE("/regproviders/:id", app.RegisterProviderHandler.DeleteRegProvider)

	r.GET("/ecr/images", app.ECRHandler.GetECRImages)

	r.GET("/dockerhub/images", app.DockerHubHandler.GetDockerHubImages)

	r.GET("/project/:id/deploy", app.ReverseProxyHandler.Forward())
	r.POST("/project/:id/build", app.BuildHandler.Build)
	r.POST("/project/:id/deploy", app.DeployHandler.Deploy)
	r.DELETE("/project/:id/deploy", app.ReverseProxyHandler.Forward())
	r.GET("/project/:id/deploy/log", app.DeployHandler.GetDeploymentLog)
	r.GET("/project/:id/deploy/:serviceName", app.ReverseProxyHandler.Forward())

	r.GET("/project/:id/deployenv", app.ReverseProxyHandler.Forward())
	r.POST("/project/:id/deployenv", app.ReverseProxyHandler.Forward())
	r.DELETE("/project/:id/deployenv", app.ReverseProxyHandler.Forward())

	r.GET("/setting/workerpool", app.ReverseProxyHandler.Forward())
	r.POST("/setting/workerpool", app.ReverseProxyHandler.Forward())

	r.GET("/github/userrepos", app.GithubAPIHandler.GetUserRepos)
	r.GET("/github/:owner/:repo/contents", app.GithubAPIHandler.GetRepoContents) // ?path={folderPath}&branch={branchName}
	r.GET("/github/:owner/:repo/branches", app.GithubAPIHandler.GetRepoBranches)
	r.POST("/github/:owner/:repo/create-branch", app.GithubAPIHandler.CreateBranch) // ?branch_name={new-branch}&selected_branch={main}
	r.PUT("/github/:owner/:repo/push", app.GithubAPIHandler.UpdateFileContent)
	r.GET("/github/:owner/:repo/commit-metadata", app.GithubAPIHandler.GetCommitMetadata)           // ?path={filePath}&branch={branchName}
	r.GET("/github/:owner/:repo/file-content", app.GithubAPIHandler.FetchFileContent)               // ?path={filePath}&branch={branchName}
	r.POST("/github/create-repo/:project_space_id/resource", app.GithubAPIHandler.CreateRepository) // might need some changes
	r.GET("/github/login", app.GithubAPIHandler.RedirectToGitHub)
	r.GET("/github/callback", app.GithubAPIHandler.GitHubCallback) // ?code={code got from the above api in search bar}

	r.GET("/project/:id/services", app.GithubAPIHandler.GetServices)

}
