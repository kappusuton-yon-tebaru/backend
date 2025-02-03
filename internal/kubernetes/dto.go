package kubernetes

type BuildImageDTO struct {
	Id           string
	Dockerfile   string
	Url          string
	Destinations []string
}
