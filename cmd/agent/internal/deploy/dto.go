package deploy

type ModifyDeploymentEnvRequest struct {
	ProjectId string
	Name      string `json:"name" validate:"required,kebabnum"`
}

type DeleteDeploymentRequest struct {
	ProjectId     string `json:"project_id" validate:"required"`
	DeploymentEnv string `json:"deployment_env" validate:"omitempty"`
	ServiceName   string `json:"service_name"   validate:"required,gt=0"`
}
