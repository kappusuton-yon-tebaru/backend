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
}
