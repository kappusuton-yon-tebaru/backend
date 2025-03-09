package deploy

type DeployRequest struct {
	ProjectId string        `json:"project_id" validate:"required"`
	Services  []ServiceInfo `json:"services"   validate:"required,gt=0,dive,required"`
}

type ServiceInfo struct {
	ServiceName string `json:"service_name" validate:"required"`
	Tag         string `json:"tag" validate:"required,ascii"`
}
