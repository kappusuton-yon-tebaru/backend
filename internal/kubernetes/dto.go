package kubernetes

type BuildImageDTO struct {
	Dockerfile   string
	Url          string
	Destinations []string
	AppName      string
}
