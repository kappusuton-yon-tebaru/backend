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
	repo *Repository
	resourceRelaRepo *resourcerelationship.Repository
}

func NewService(repo *Repository,resourceRelaRepo *resourcerelationship.Repository) *Service {
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
	
	resource, err := s.repo.GetResourceByID(ctx, filter)
	if err != nil {
		return models.Resource{}, werror.NewFromError(err)
	}

	return resource, nil
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
	// get children rela
	childrenResourceRelas, err := s.resourceRelaRepo.GetChildrenResourceRelationshipByParentID(ctx,filter)
	if err != nil {
		return nil, werror.NewFromError(err)
	}
	childrenResources := []models.Resource{}
	//get children resource
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
		childrenResource, err := s.repo.GetResourceByID(ctx,filter)
		if err != nil {
			return nil, werror.NewFromError(err)
		}
		childrenResources = append(childrenResources, childrenResource )
	}

	return childrenResources, nil
}

func (s *Service) CreateResource(ctx context.Context, dto CreateResourceDTO, id string) (string, error) {
	resource_id, err := s.repo.CreateResource(ctx, dto)
	if err != nil {
		return "", err
	}
	// Is an org no need to create rela
	if id == "" {
		return resource_id, nil
	}
	
	parentID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Invalid parent ID:", err)
		return  "", err
	}

	childID, err := bson.ObjectIDFromHex(resource_id)
	if err != nil {
		fmt.Println("Invalid child ID:", err)
		return  "", err
	}

	// Create DTO instance
	relationship := resourcerelationship.CreateResourceRelationshipDTO{
		Parent_Resource_Id: parentID,
		Child_Resource_Id:  childID,
	}

	resource_rela_id, err := s.resourceRelaRepo.CreateResourceRelationship(ctx,relationship)
	if err != nil {
		return "", err
	}
	
	return resource_rela_id, nil
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
