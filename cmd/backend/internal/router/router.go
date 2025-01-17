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
}
