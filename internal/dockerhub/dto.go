package dockerhub

type GetDockerHubImagesRequest struct {
	Namespace      string `json:"namespace"`
	RepositoryName string `json:"repository_name"`
	ServiceName    string `json:"service_name"`
}

type DockerHubImageResponse struct {
	ImageTag    string `json:"image_tag"`
	LastUpdated string `json:"last_updated"`
}
