package githubapi

import (
	"context"
	"errors"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) GetUserRepos(ctx context.Context, token string) ([]models.Repository, error) {
	return s.repo.GetUserRepos(ctx, token)
}

func (s *Service) GetRepoContents(ctx context.Context, fullname string, path string, token string) ([]models.File, error) {
    if fullname == "" {
        return nil, errors.New("Repository fullname is required")
    }

    return s.repo.GetRepoContents(fullname, path, token)
}

// GetRepoBranches fetches the branches of a repository
func (s *Service) GetRepoBranches(ctx context.Context, fullname string, token string) ([]models.Branch, error) {
    if fullname == "" {
        return nil, errors.New("Repository fullname is required")
    }

    return s.repo.GetRepoBranches(fullname, token)
}

// GetCommitMetadata fetches the commit metadata for a file in a repository
func (s *Service) GetCommitMetadata(ctx context.Context, path string, branch string, fullname string, token string) (*models.CommitMetadata, error) {
    if fullname == "" || path == "" || branch == "" {
        return nil, errors.New("Repository fullname, path, and branch are required")
    }

    return s.repo.GetCommitMetadata(path, branch, fullname, token)
}

func (s *Service) FetchFileContent(ctx context.Context, fullname, filePath, branch, token string) (string, string, error) {
	if fullname == "" || filePath == "" || branch == "" || token == "" {
		return "", "", errors.New("missing required parameters")
	}

	return s.repo.FetchFileContent(ctx, fullname, filePath, branch, token)
}