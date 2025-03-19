package deploy

type DeployRequest struct {
	ProjectId     string        `swaggerignore:"true"`
	DeploymentEnv string        `json:"deployment_env" validate:"omitempty"`
	Services      []ServiceInfo `json:"services"   validate:"required,gt=0,dive,required"`
}

type ServiceInfo struct {
	ServiceName string  `json:"service_name" validate:"required"`
	Tag         string  `json:"tag" validate:"required,ascii"`
	Port        *int32  `json:"port" validate:"omitempty,min=1"`
	SecretName  *string `json:"secret_name" validate:"omitempty,ascii"`
}

type DeployResponse struct {
	ParentId string `json:"parent_id"`
}
