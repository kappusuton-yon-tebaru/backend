package resource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/resourcerelationship"

	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PaginatedResources = models.Paginated[models.Resource]
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

func (s *Service) GetChildrenResourcesByParentID(ctx context.Context, queryParam query.QueryParam, parentId string) (PaginatedResources, error) {
	dtos, err := s.repo.GetResourcesByFilter(ctx, queryParam, parentId)
	if err != nil {
		return PaginatedResources{}, err
	}
	resources := make([]models.Resource, 0)
	for _, dto := range dtos.Data {
		fmt.Println(DTOToResource(dto))
		resources = append(resources, DTOToResource(dto))
	}

	return PaginatedResources{
		Page:  dtos.Page,
		Limit: dtos.Limit,
		Total: dtos.Total,
		Data:  resources,
	}, nil
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
