package deploy

type DeployRequest struct {
	ProjectId string        `json:"project_id" validate:"required"`
	Services  []ServiceInfo `json:"services"   validate:"required,gt=0,dive,required"`
}

type ServiceInfo struct {
	ServiceName string  `json:"service_name" validate:"required"`
	Tag         string  `json:"tag" validate:"required,ascii"`
	Port        *int    `json:"port" validate:"omitempty,min=1"`
	SecretName  *string `json:"secret_name" validate:"omitempty,ascii"`
}
