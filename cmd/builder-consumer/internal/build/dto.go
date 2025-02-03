package build

type BuildRequestDTO struct {
	Id           string   `json:"id"`
	Dockerfile   string   `json:"dockerfile"`
	Url          string   `json:"url"`
	Destinations []string `json:"destinations"`
}
