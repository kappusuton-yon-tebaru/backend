package build

type BuildRequest struct {
	ProjectId string        `json:"project_id" validate:"required"`
	Services  []ServiceInfo `json:"services"   validate:"required,gt=0,dive,required"`
}

type BuildResponse struct {
	ParentId string `json:"parent_id"`
}

type ServiceInfo struct {
	ServiceName string `json:"service_name" validate:"required,filepath"`
	Tag         string `json:"tag"          validate:"required,ascii"`
}
