package githubapi

import (
	"context"

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
