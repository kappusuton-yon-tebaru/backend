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
	r.DELETE("/resources/:id", app.ResourceHandler.DeleteResource)

	r.GET("/roles", app.RoleHandler.GetAllRoles)
	r.DELETE("/roles/:id", app.RoleHandler.DeleteRoleById)

	r.GET("projrepos", app.ProjectRepositoryHandler.GetAllProjectRepositories)
	r.DELETE("projrepos/:id", app.ProjectRepositoryHandler.DeleteProjectRepository)

	r.GET("resourcerelas", app.ResourceRelationshipHandler.GetAllResourceRelationships)
	r.DELETE("resourcerelas/:id", app.ResourceRelationshipHandler.DeleteResourceRelationship)
}
