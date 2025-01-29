package projectenv

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

func (s *Service) GetAllProjectEnvs(ctx context.Context) ([]models.ProjectEnv, error) {
	projectenvs, err := s.repo.GetAllProjectEnvs(ctx)
	if err != nil {
		return nil, err
	}

	return projectenvs, nil
}

func (s *Service) CreateProjectEnv(ctx context.Context, dto CreateProjectEnvDTO) (any, error) {
	id, err := s.repo.CreateProjectEnv(ctx, dto)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) DeleteProjectEnv(ctx context.Context, id string) *werror.WError {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid projectenv id")
	}

	filter := map[string]any{
		"_id": objId,
	}

	count, err := s.repo.DeleteProjectEnv(ctx, filter)
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