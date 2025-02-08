package kubernetes

type BuildImageDTO struct {
	Id           string
	Dockerfile   string
	RepoUrl      string
	Destinations []string
}
