package kubernetes

import "github.com/kappusuton-yon-tebaru/backend/internal/models"

type BuildImageDTO struct {
	Id            string
	Dockerfile    string
	RepoUrl       string
	RepoRoot      string
	Destinations  []string
	ECRCredential *models.ECRCredential
}

type DeployDTO struct {
	ServiceName  string
	ImageUri     string
	Port         *int32
	Namespace    string
	Environments map[string]string
}

type ConfigureMaxWorkerDTO struct {
	WorkerImageUri string
	WorkerNumber   int32
}
