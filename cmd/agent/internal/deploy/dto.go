package deploy

type ModifyDeploymentEnvRequest struct {
	ProjectId string `swaggerignore:"true"`
	Name      string `json:"name" validate:"required,kebabnum"`
}

type DeleteDeploymentRequest struct {
	ProjectId     string `swaggerignore:"true"`
	DeploymentEnv string `json:"deployment_env" validate:"omitempty"`
	ServiceName   string `json:"service_name"   validate:"required,gt=0"`
}

type ListDeploymentEnvResponse struct {
	Data []string `json:"data"`
}
