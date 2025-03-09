package ecr

type GetECRImagesRequest struct {
	ProjectId   string `json:"project_id"`
	ServiceName string `json:"service_name"`
}

type ECRImageResponse struct {
	ImageTag string `json:"image_tag"`
}
