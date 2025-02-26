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

func (s *Service) GetChildrenResourceRelationshipByParentID(
	ctx context.Context, parentID string, page, limit int,
) ([]models.ResourceRelationship, int, *werror.WError) {

	objId, err := bson.ObjectIDFromHex(parentID)
	if err != nil {
		return nil, 0, werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid parent id")
	}

	filter := map[string]any{
		"parent_resource_id": objId,
	}

	offset := (page - 1) * limit
	childrenResourceRelas, total, err := s.repo.GetChildrenResourceRelationshipByParentID(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, werror.NewFromError(err)
	}

	return childrenResourceRelas, total, nil
}

func (s *Service) CreateResourceRelationship(ctx context.Context, dto CreateResourceRelationshipDTO) (string, error) {
	id, err := s.repo.CreateResourceRelationship(ctx, dto)
	if err != nil {
		return "", err
	}

	return id, nil
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
