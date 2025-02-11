package resource

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

func (s *Service) GetAllResources(ctx context.Context) ([]models.Resource, error) {
	resources, err := s.repo.GetAllResources(ctx)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (s *Service) GetResourceByID(ctx context.Context, id string) (models.Resource, *werror.WError) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Resource{}, werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid id")
	}

	filter := map[string]any{
		"_id": objId,
	}
	
	resource, err := s.repo.GetResourceByID(ctx, filter)
	if err != nil {
		return models.Resource{}, werror.NewFromError(err)
	}

	return resource, nil
}

func (s *Service) CreateResource(ctx context.Context, dto CreateResourceDTO) (string, error) {
	id, err := s.repo.CreateResource(ctx, dto)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) DeleteResource(ctx context.Context, id string) *werror.WError {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid resource id")
	}

	filter := map[string]any{
		"_id": objId,
	}

	count, err := s.repo.DeleteResource(ctx, filter)
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
