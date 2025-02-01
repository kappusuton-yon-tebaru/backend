package build

type BuildRequestDTO struct {
	AppName      string   `json:"app_name"`
	Dockerfile   string   `json:"dockerfile"`
	Url          string   `json:"url"`
	Destinations []string `json:"destinations"`
}
