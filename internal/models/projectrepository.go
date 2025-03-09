package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type ProjectRepository struct {
	Id               string             `json:"id"`
	ProjectId        string             `json:"project_id"`
	GitRepoUrl       string             `json:"git_repo_url"`
	RegistryProvider *RegistryProviders `json:"registry_provider"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

func (pr ProjectRepository) GetGitRepoUrl() (string, error) {
	rg := regexp.MustCompile("github.com/(?<owner>[^/]+)/(?<repo>[^/]+)$")

	matches := rg.FindStringSubmatch(pr.GitRepoUrl)

	if len(matches) != 3 {
		return "", errors.New("error occured while parsing git repo url")
	}

	owner := matches[rg.SubexpIndex("owner")]
	repo := strings.TrimSuffix(matches[rg.SubexpIndex("repo")], ".git")

	return fmt.Sprintf("git://github.com/%s/%s", owner, repo), nil
}
