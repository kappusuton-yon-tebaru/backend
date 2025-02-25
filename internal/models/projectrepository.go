package models

import "time"

type ProjectRepository struct {
	Id               string            `json:"id"`
	ProjectId        string            `json:"project_id"`
	GitRepoUrl       string            `json:"git_repo_url"`
	RegistryProvider RegistryProviders `json:"registry_provider"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}
