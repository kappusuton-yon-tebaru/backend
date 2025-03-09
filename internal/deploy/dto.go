package build

type DeployContext struct {
	Id           string            `json:"id"`
	ServiceName  string            `json:"service_name"`
	ImageUri     string            `json:"image_uri"`
	Environments map[string]string `json:"environments"`
}
