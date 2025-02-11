package resourcerelationship

import (
	"context"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/resource"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	repo *Repository
	resourceRepo *resource.Repository
}

func NewService(repo *Repository, resourceRepo *resource.Repository) *Service {
	return &Service{
		repo,
		resourceRepo,
	}
}

func (s *Service) GetAllResourceRelationships(ctx context.Context) ([]models.ResourceRelationship, error) {
	resourceRelas, err := s.repo.GetAllResourceRelationships(ctx)
	if err != nil {
		return nil, err
	}

	return resourceRelas, nil
}

func (s *Service) GetChildrenResourcesByParentID(ctx context.Context, parentID string) ([]models.Resource, *werror.WError) {
	objId, err := bson.ObjectIDFromHex(parentID)
	if err != nil {
		return nil, werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid parent id")
	}

	filter := map[string]any{
		"parent_resource_id": objId,
	}
	
	childrenResourceRelas, err := s.repo.GetChildrenResourcesByParentID(ctx, filter)
	if err != nil {
		return nil, werror.NewFromError(err)
	}

	var childrenResources []models.Resource

	for _, childrenResourceRela := range childrenResourceRelas{
		objId, err := bson.ObjectIDFromHex(childrenResourceRela.ChildResourceId)
		if err != nil {
			return nil, werror.NewFromError(err).
				SetCode(http.StatusBadRequest).
				SetMessage("invalid id")
		}

		filter := map[string]any{
			"_id": objId,
		}
		childrenResource, err := s.resourceRepo.GetResourceByID(ctx,filter)
		if err != nil {
			return nil, werror.NewFromError(err)
		}
		childrenResources = append(childrenResources, childrenResource )
	}

	return childrenResources, nil
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
