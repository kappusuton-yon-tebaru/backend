package models

type ProjectRepository struct {
	Id         string `json:"id"`
	GitRepoUrl string `json:"git_repo_url"`
	ProjectId  string `json:"project_id"`
}
