package deploy

type DeleteDeploymentRequest struct {
	ProjectId     string `swaggerignore:"true"`
	DeploymentEnv string `json:"deployment_env" validate:"omitempty"`
	ServiceName   string `json:"service_name"   validate:"required,gt=0"`
}

type ListDeploymentQuery struct {
	ProjectId     string
	DeploymentEnv string `form:"deployment_env"`
}

type DeploymentResponse struct {
	Message string `json:"message"`
}

type GetServiceDeployment struct {
	ProjectId     string
	DeploymentEnv string `form:"deployment_env"`
	ServiceName   string
}
