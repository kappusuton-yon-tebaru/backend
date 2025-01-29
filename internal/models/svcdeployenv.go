package models

type ServiceDeploymentEnv struct {
	Id                  string `json:"id"`
	ServiceDeploymentId string `json:"service_deployment_id"`
	EnvType             string `json:"env_type"`
	Key                 string `json:"key"`
	Value               string `json:"value"`
}
