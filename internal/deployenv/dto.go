package deployenv

type ModifyDeploymentEnvDTO struct {
	ProjectId string
	Name      string
}

type DeleteDeploymentDTO struct {
	ProjectId     string
	DeploymentEnv string
	ServiceName   string
}
