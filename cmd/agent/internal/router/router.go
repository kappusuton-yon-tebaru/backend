package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/cmd/agent/agent"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
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

func (r *Router) RegisterRoutes(app *agent.App) {
	r.GET("/ws/job/:id/log", app.MonitoringHandler.IntervalPing)
	// r.GET("/ws/job/:id/log", app.MonitoringHandler.StreamJobLog)
}
