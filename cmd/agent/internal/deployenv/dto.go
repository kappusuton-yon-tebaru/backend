package deployenv

type ModifyDeploymentEnvRequest struct {
	ProjectId string `swaggerignore:"true"`
	Name      string `json:"name" validate:"required,kebabnum"`
}

type ListDeploymentEnvResponse struct {
	Data []string `json:"data"`
}

type DeploymentDevResponse struct {
	Message string `json:"message"`
}
