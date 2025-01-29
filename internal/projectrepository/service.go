package projectrepository

import (
	"context"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) GetAllProjectRepositories(ctx context.Context) ([]models.ProjectRepository, error) {
	projRepos, err := s.repo.GetAllProjectRepositories(ctx)
	if err != nil {
		return nil, err
	}

	return projRepos, nil
}

func (s *Service) CreateProjectRepository(ctx context.Context, dto CreateProjectRepositoryDTO) (any, error) {
	id, err := s.repo.CreateProjectRepository(ctx, dto)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) DeleteProjectRepository(ctx context.Context, id string) *werror.WError {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid project repository id")
	}

	filter := map[string]any{
		"_id": objId,
	}

	count, err := s.repo.DeleteProjectRepository(ctx, filter)
	if err != nil {
		return werror.NewFromError(err)
	}

	if count == 0 {
		return werror.New().
			SetCode(http.StatusNotFound).
			SetMessage("not found")
	}

	return nil
}
