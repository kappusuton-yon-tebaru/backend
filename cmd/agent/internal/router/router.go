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
	r.GET("/setting/workerpool", app.SettingHandler.GetWorkerPoolSetting)
	r.POST("/setting/workerpool", app.SettingHandler.SetWorkerPoolSetting)

	r.GET("/project/:id/deployenv", app.DeployEnvHandler.ListDeploymentEnv)
	r.POST("/project/:id/deployenv", app.DeployEnvHandler.CreateDeploymentEnv)
	r.DELETE("/project/:id/deployenv", app.DeployEnvHandler.DeleteDeploymentEnv)

	r.GET("/project/:id/deploy/:serviceName", app.DeployHandler.GetServiceDeployment)
	r.GET("/project/:id/deploy", app.DeployHandler.ListDeployment)
	r.DELETE("/project/:id/deploy", app.DeployHandler.DeleteDeployment)
}
