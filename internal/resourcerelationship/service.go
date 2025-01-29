package resourcerelationship

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

func (s *Service) GetAllResourceRelationships(ctx context.Context) ([]models.ResourceRelationship, error) {
	resourceRelas, err := s.repo.GetAllResourceRelationships(ctx)
	if err != nil {
		return nil, err
	}

	return resourceRelas, nil
}

func (s *Service) DeleteResourceRelationship(ctx context.Context, id string) *werror.WError {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid resource relationship id")
	}

	filter := map[string]any{
		"_id": objId,
	}

	count, err := s.repo.DeleteResourceRelationship(ctx, filter)
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
