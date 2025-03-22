package models

import (
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
)

type Deployment struct {
	ProjectId     string                `json:"project_id"`
	ProjectName   string                `json:"project_name"`
	ServiceName   string                `json:"service_name"`
	DeploymentEnv string                `json:"deployment_env"`
	Status        enum.DeploymentStatus `json:"deployment_status"`
	Age           time.Duration         `json:"-"`
	StringAge     string                `json:"age"`
	// ServiceUrl string `json:"service_url"`
}
