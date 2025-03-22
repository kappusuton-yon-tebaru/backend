package deploy

type DeployContext struct {
	Id            string                    `json:"id"`
	ProjectId     string                    `json:"project_id"`
	ServiceName   string                    `json:"service_name"`
	ImageUri      string                    `json:"image_uri"`
	Port          *int32                    `json:"port,omitempty"`
	Namespace     string                    `json:"namespace"`
	DeploymentEnv string                    `json:"deploy_env"`
	Environments  map[string]string         `json:"environments"`
	HealthCheck   *DeployHealthCheckContext `json:"health_check,omitempty"`
}

type DeployHealthCheckContext struct {
	Path string `json:"path"`
	Port int32  `json:"port"`
}
