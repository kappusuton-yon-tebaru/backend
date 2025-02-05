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

func (s *Service) GetRepoContents(fullname string, path string, token string) ([]models.File, error) {
    if fullname == "" {
        return nil, errors.New("Repository fullname is required")
    }

    return s.repo.GetRepoContents(fullname, path, token)
}
