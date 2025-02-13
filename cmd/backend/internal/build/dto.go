package build

type BuildRequest struct {
	RepoUrl     string        `json:"repo_url"     validate:"url,required"`
	RegistryUrl string        `json:"registry_url" validate:"required"`
	Services    []ServiceInfo `json:"services"     validate:"required,gt=0,dive,required"`
}

type ServiceInfo struct {
	Dockerfile string `json:"dockerfile"   validate:"required,filepath"`
	Tag        string `json:"tag"          validate:"required,ascii"`
}
