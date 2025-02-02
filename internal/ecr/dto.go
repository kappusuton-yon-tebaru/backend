package ecr

type PushECRImageRequest struct {
	RepositoryName string `json:"repository_name"`
	ImageName string `json:"image_name"`
	Tag       string `json:"tag"`
}

type GetECRImagesRequest struct {
	RepositoryName string `json:"repository_name"`
}

type ECRImageResponse struct {
	RepositoryName string `json:"repository_name"`
	ImageDigest    string `json:"image_digest"`
	ImageTag       string `json:"image_tag"`
}
