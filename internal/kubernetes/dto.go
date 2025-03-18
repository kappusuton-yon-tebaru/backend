package kubernetes

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"k8s.io/api/apps/v1"
)

type BuildImageDTO struct {
	Id            string
	Dockerfile    string
	RepoUrl       string
	RepoRoot      string
	Destinations  []string
	ECRCredential *models.ECRCredential
}

type DeployDTO struct {
	Id           string
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

type DeploymentCondition struct {
	Available      *v1.DeploymentCondition
	Progressing    *v1.DeploymentCondition
	ReplicaFailure *v1.DeploymentCondition
}
