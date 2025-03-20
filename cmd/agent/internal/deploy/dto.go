package deploy

type DeleteDeploymentRequest struct {
	ProjectId     string `swaggerignore:"true"`
	DeploymentEnv string `json:"deployment_env" validate:"omitempty"`
	ServiceName   string `json:"service_name"   validate:"required,gt=0"`
}
