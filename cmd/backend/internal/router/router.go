package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/backend"
)

type Router struct {
	*gin.Engine
}

func New() *Router {
	return &Router{
		gin.Default(),
	}
}

func (r *Router) RegisterRoutes(app *backend.App) {
	r.GET("/", app.GreetingHandler.Greeting)

	r.GET("/images", app.ImageHandler.GetAllImages)
	r.DELETE("/image/:id", app.ImageHandler.DeleteImage)

	r.GET("/svcdeploys", app.ServiceDeployment.GetAllServiceDeployments)
	r.DELETE("/svcdeploy/:id", app.ServiceDeployment.DeleteServiceDeployment)

	r.GET("/users", app.UserHandler.GetAllUsers)
	r.POST("/users", app.UserHandler.CreateUser)
	r.DELETE("/users/:id", app.UserHandler.DeleteUserById)

	r.GET("/usergroups", app.UserGroupHandler.GetAllUserGroups)
	r.POST("/usergroups", app.UserGroupHandler.CreateUserGroup)
	r.DELETE("/usergroups/:group_id", app.UserGroupHandler.DeleteUserGroupById)
	r.POST("/usergroups/:id/user", app.UserGroupHandler.AddUserToUserGroup)
	r.DELETE("/usergroups/:group_id/user/:user_id", app.UserGroupHandler.DeleteUserFromUserGroupById)

	r.GET("/resources", app.ResourceHandler.GetAllResources)
	r.POST("/resources", app.ResourceHandler.CreateResource)
	r.DELETE("/resources/:id", app.ResourceHandler.DeleteResource)

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
	r.POST("/projrepos", app.ProjectRepositoryHandler.CreateProjectRepository)
	r.DELETE("/projrepos/:id", app.ProjectRepositoryHandler.DeleteProjectRepository)

	r.GET("/resourcerelas", app.ResourceRelationshipHandler.GetAllResourceRelationships)
	r.POST("/resourcerelas", app.ResourceRelationshipHandler.CreateResourceRelationship)
	r.DELETE("/resourcerelas/:id", app.ResourceRelationshipHandler.DeleteResourceRelationship)

	r.GET("/jobs", app.JobHandler.GetAllJobs)
	r.POST("/jobs", app.JobHandler.CreateJob)
	r.DELETE("/jobs/:id", app.JobHandler.DeleteJob)

	r.GET("/regproviders", app.RegisterProviderHandler.GetAllRegProviders)
	r.POST("/regproviders", app.RegisterProviderHandler.CreateRegProvider)
	r.DELETE("/regproviders/:id", app.RegisterProviderHandler.DeleteRegProvider)

	r.GET("/projectenvs", app.ProjectEnvironmentHandler.GetAllProjectEnvs)
	r.POST("/projectenvs", app.ProjectEnvironmentHandler.CreateProjectEnv)
	r.DELETE("/projectenvs/:id", app.ProjectEnvironmentHandler.DeleteProjectEnv)
}
