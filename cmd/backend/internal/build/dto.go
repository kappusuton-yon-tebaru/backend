package build

// type BuildRequest struct {
// 	RepoUrl     string        `json:"repo_url"     validate:"url,required"`
// 	RegistryUrl string        `json:"registry_url" validate:"required"`
// 	Services    []ServiceInfo `json:"services"     validate:"required,gt=0,dive,required"`
// }

type BuildRequest struct {
	ProjectId string        `json:"project_id" validate:"required"`
	Services  []ServiceInfo `json:"services"   validate:"required,gt=0,dive,required"`
}

// type ServiceInfo struct {
// 	ServiceRoot string `json:"service_root" validate:"required,filepath"`
// 	Tag         string `json:"tag"          validate:"required,ascii"`
// }

type ServiceInfo struct {
	ServiceName string `json:"service_name" validate:"required,filepath"`
	Tag         string `json:"tag"          validate:"required,ascii"`
}
