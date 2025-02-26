package resource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/resourcerelationship"

	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	repo             *Repository
	resourceRelaRepo *resourcerelationship.Repository
}

func NewService(repo *Repository, resourceRelaRepo *resourcerelationship.Repository) *Service {
	return &Service{
		repo,
		resourceRelaRepo,
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

	resource, err := s.repo.GetResourceByFilter(ctx, filter)
	if err != nil {
		return models.Resource{}, werror.NewFromError(err)
	}

	return resource, nil
}

func (s *Service) GetChildrenResourcesByParentID(ctx context.Context, parentID string, page, limit int) ([]models.Resource, int, *werror.WError) {
	objId, err := bson.ObjectIDFromHex(parentID)
	if err != nil {
		return nil, 0, werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid parent id")
	}

	filter := map[string]any{
		"parent_resource_id": objId,
	}

	skip := (page - 1) * limit

	childrenResourceRelas, total, err := s.resourceRelaRepo.GetChildrenResourceRelationshipByParentID(ctx, filter, limit, skip)
	if err != nil {
		return nil, 0, werror.NewFromError(err)
	}

	childrenResources := []models.Resource{}

	for _, childrenResourceRela := range childrenResourceRelas {
		objId, err := bson.ObjectIDFromHex(childrenResourceRela.ChildResourceId)
		if err != nil {
			return nil, 0, werror.NewFromError(err).
				SetCode(http.StatusBadRequest).
				SetMessage("invalid id")
		}

		filter := map[string]any{
			"_id": objId,
		}
		childrenResource, err := s.repo.GetResourceByFilter(ctx, filter)
		if err != nil {
			return nil, 0, werror.NewFromError(err)
		}
		childrenResources = append(childrenResources, childrenResource)
	}

	return childrenResources, total, nil
}

func (s *Service) CreateResource(ctx context.Context, dto CreateResourceDTO, parentID string) (string, error) {
	resourceId, err := s.repo.CreateResource(ctx, dto)
	if err != nil {
		return "", err
	}
	// Is an org no need to create rela
	if parentID == "" {
		return resourceId, nil
	}

	pID, err := bson.ObjectIDFromHex(parentID)
	if err != nil {
		fmt.Println("Invalid parent ID:", err)
		return "", err
	}

	childID, err := bson.ObjectIDFromHex(resourceId)
	if err != nil {
		fmt.Println("Invalid child ID:", err)
		return "", err
	}

	// Create DTO instance
	relationship := resourcerelationship.CreateResourceRelationshipDTO{
		ParentResourceId: pID,
		ChildResourceId:  childID,
	}

	_, err = s.resourceRelaRepo.CreateResourceRelationship(ctx, relationship)
	if err != nil {
		return "", err
	}

	return resourceId, nil
}

func (s *Service) UpdateResource(ctx context.Context, dto UpdateResourceDTO, id string) (string, *werror.WError) {
	resourceId, err := s.repo.UpdateResource(ctx, dto, id)
	if err != nil {
		return "", werror.NewFromError(err).
			SetCode(http.StatusBadRequest)
	}

	return resourceId, nil
}

func (s *Service) DeleteResource(ctx context.Context, id string) *werror.WError {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid resource id")
	}

	deleteFilter := map[string]any{
		"_id": objId,
	}

	count, err := s.repo.DeleteResource(ctx, deleteFilter)
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

func (s *Service) CascadeDeleteResource(ctx context.Context, id string) *werror.WError {
	resource, err := s.GetResourceByID(ctx, id)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("resource not found")
	}

	err2 := s.repo.CascadeDeleteResource(ctx, id, resource.ResourceType)
	if err2 != nil {
		return werror.NewFromError(err2)
	}

	return nil
}
