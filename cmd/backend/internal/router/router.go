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

	r.GET("/", app.GreetingHandler.Greeting)

	r.GET("/images", app.ImageHandler.GetAllImages)
	r.DELETE("/image/:id", app.ImageHandler.DeleteImage)

	r.GET("/svcdeploys", app.ServiceDeployment.GetAllServiceDeployments)
	r.DELETE("/svcdeploy/:id", app.ServiceDeployment.DeleteServiceDeployment)

	r.GET("/users", app.UserHandler.GetAllUsers)
	r.DELETE("/users/:id", app.UserHandler.DeleteUserById)

	r.POST("/register", app.UserHandler.Register)

	r.GET("/usergroups", app.UserGroupHandler.GetAllUserGroups)
	r.POST("/usergroups", app.UserGroupHandler.CreateUserGroup)
	r.DELETE("/usergroups/:group_id", app.UserGroupHandler.DeleteUserGroupById)
	r.POST("/usergroups/:id/user", app.UserGroupHandler.AddUserToUserGroup)
	r.DELETE("/usergroups/:group_id/user/:user_id", app.UserGroupHandler.DeleteUserFromUserGroupById)

	r.GET("/resources", app.ResourceHandler.GetAllResources)
	r.GET("/resources/:id", app.ResourceHandler.GetResourceByID)
	r.GET("/resources/children/:parent_id", app.ResourceHandler.GetChildrenResourcesByParentID)
	r.POST("/resources", app.ResourceHandler.CreateResource) // ?parent_id={id}
	r.PUT("/resources/:id", app.ResourceHandler.UpdateResource)
	r.DELETE("/resources/:id", app.ResourceHandler.DeleteResource)
	r.DELETE("/resources/cascade/:id", app.ResourceHandler.CascadeDeleteResource)

	r.GET("/roles", app.RoleHandler.GetAllRoles)
	r.POST("/roles", app.RoleHandler.CreateRole)
	r.DELETE("/roles/:id", app.RoleHandler.DeleteRoleById)

	r.GET("/permissions", app.PermissionHandler.GetAllPermissions)
	r.POST("/permissions", app.PermissionHandler.CreatePermission)
	r.DELETE("/permissions/:id", app.PermissionHandler.DeletePermissionById)

	r.GET("/rolepermissions", app.RolePermissionHandler.GetAllRolePermissions)
	r.POST("/rolepermissions", app.RolePermissionHandler.CreateRolePermission)
	r.DELETE("/rolepermissions/:id", app.RolePermissionHandler.DeleteRolePermissionById)

	r.GET("/roleusergroups", app.RoleUserGroupHandler.GetAllRoleUserGroups)
	r.POST("/roleusergroups", app.RoleUserGroupHandler.CreateRoleUserGroup)
	r.DELETE("/roleusergroups/:id", app.RoleUserGroupHandler.DeleteRoleUserGroupById)

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
	r.GET("/jobs/:id", app.JobHandler.GetAllJobsByParentId)

	r.GET("/regproviders", app.RegisterProviderHandler.GetAllRegProviders)
	r.GET("/regproviders/unbind", app.RegisterProviderHandler.GetAllRegProvidersWithoutProject)
	r.GET("/regproviders/:id", app.RegisterProviderHandler.GetRegProviderById)
	r.POST("/regproviders", app.RegisterProviderHandler.CreateRegProvider)
	r.DELETE("/regproviders/:id", app.RegisterProviderHandler.DeleteRegProvider)

	r.GET("/projectenvs", app.ProjectEnvironmentHandler.GetAllProjectEnvs)
	r.POST("/projectenvs", app.ProjectEnvironmentHandler.CreateProjectEnv)
	r.DELETE("/projectenvs/:id", app.ProjectEnvironmentHandler.DeleteProjectEnv)

	r.GET("/ecr/images", app.ECRHandler.GetECRImages)

	r.GET("/dockerhub/images", app.DockerHubHandler.GetDockerHubImages)

	r.GET("/project/:id/deploy", app.ReverseProxyHandler.Forward())
	r.POST("/project/:id/build", app.BuildHandler.Build)
	r.POST("/project/:id/deploy", app.DeployHandler.Deploy)
	r.DELETE("/project/:id/deploy", app.ReverseProxyHandler.Forward())

	r.GET("/project/:id/deployenv", app.ReverseProxyHandler.Forward())
	r.POST("/project/:id/deployenv", app.ReverseProxyHandler.Forward())
	r.DELETE("/project/:id/deployenv", app.ReverseProxyHandler.Forward())

	r.GET("/ws/job/:id/log", app.MonitoringHandler.StreamJobLog)

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
