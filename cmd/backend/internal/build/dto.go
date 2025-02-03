package build

type BuildRequest struct {
	Id           string   `json:"id"`
	Dockerfile   string   `json:"dockerfile"   validate:"required"`
	Url          string   `json:"url"          validate:"required"`
	Destinations []string `json:"destinations" validate:"required"`
}
