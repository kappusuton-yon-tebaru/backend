package build

type BuildContext struct {
	Id          string `json:"id"`
	RepoUrl     string `json:"repo_url"`
	Destination string `json:"destination"`
	Dockerfile  string `json:"dockerfile"`
}
