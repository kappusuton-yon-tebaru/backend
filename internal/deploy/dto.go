package build

type DeployContext struct {
	Id           string            `json:"id"`
	ServiceName  string            `json:"service_name"`
	ImageUri     string            `json:"image_uri"`
	Port         *int              `json:"port"`
	Environments map[string]string `json:"environments"`
}
