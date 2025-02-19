package ecr

type GetECRImagesRequest struct {
	// RepositoryName 	string `json:"repository_name"`
	// RepositoryURI 	string `json:"repository_uri"`
	ProjectId		string `json:"project_id"`
	ServiceName		string `json:"service_name"`
}

type ECRImageResponse struct {
	ImageTag       string `json:"image_tag"`
}
