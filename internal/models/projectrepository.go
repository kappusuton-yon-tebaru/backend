package models

import (
	"fmt"
	"regexp"
	"strings"
)

type ProjectRepository struct {
	Id               string             `json:"id"`
	ProjectId        string             `json:"project_id"`
	GitRepoUrl       string             `json:"git_repo_url"`
	RegistryProvider *RegistryProviders `json:"registry_provider"`
}

func (pr ProjectRepository) GetGitRepoUri() string {
	rg := regexp.MustCompile("github.com/(?<owner>[^/]+)/(?<repo>[^/]+)$")

	matches := rg.FindStringSubmatch(pr.GitRepoUrl)

	owner := matches[rg.SubexpIndex("owner")]
	repo := strings.TrimSuffix(matches[rg.SubexpIndex("repo")], ".git")

	return fmt.Sprintf("git://github.com/%s/%s", owner, repo)
}
