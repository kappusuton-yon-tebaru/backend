package build

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
)

type BuildContext struct {
	Id            string                `json:"id"`
	RepoUrl       string                `json:"repo_url"`
	Dockerfile    string                `json:"dockerfile"`
	RepoRoot      string                `json:"repo_root"`
	Destination   string                `json:"destination"`
	ECRCredential *models.ECRCredential `json:"ecr_credential,omitempty"`
}
