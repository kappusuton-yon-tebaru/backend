package kubernetes

type BuildImageDTO struct {
	Id           string
	Dockerfile   string
	RepoUrl      string
	RepoRoot     string
	Destinations []string
}

type ConfigureMaxWorkerDTO struct {
	WorkerImageUri string
	WorkerNumber   int32
}
