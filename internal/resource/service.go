package resource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/resourcerelationship"
	"github.com/kappusuton-yon-tebaru/backend/internal/role"
	"github.com/kappusuton-yon-tebaru/backend/internal/user"

	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PaginatedResources = models.Paginated[models.Resource]
type Service struct {
	repo             *Repository
	resourceRelaRepo *resourcerelationship.Repository
	roleRepo         *role.Repository
	userRepo         *user.Repository
}

func NewService(repo *Repository, resourceRelaRepo *resourcerelationship.Repository, roleRepo *role.Repository, userRepo *user.Repository) *Service {
	return &Service{
		repo,
		resourceRelaRepo,
		roleRepo,
		userRepo,
	}
}

func (s *Service) GetAllResources(ctx context.Context, ids []string) ([]models.Resource, error) {
	objIds := make([]bson.ObjectID, 0)
	for _, id := range ids {
		objId, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return []models.Resource{}, werror.NewFromError(err).
				SetCode(http.StatusBadRequest).
				SetMessage("invalid id")
		}
		objIds = append(objIds, objId)
	}
	filter := map[string]any{"_id": map[string]any{"$in": objIds}}
	resources, err := s.repo.GetAllResources(ctx,filter)
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

func (s *Service) GetChildrenResourcesByParentID(ctx context.Context, queryParam query.QueryParam, parentId string, ids []string) (PaginatedResources, error) {
	dtos, err := s.repo.GetResourcesByFilter(ctx, queryParam, parentId,ids)
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
// use this only when creating org 
func (s *Service) AddDefaultRoleToUser(ctx context.Context, userId string, orgId string, resourceId string, name string) (error) {
	orgID, err := bson.ObjectIDFromHex(orgId)
	if err != nil {
		fmt.Println("Invalid org ID:", err)
		return err
	}
	resourceID, err := bson.ObjectIDFromHex(resourceId)
	if err != nil {
		fmt.Println("Invalid resource ID:", err)
		return err
	}
	// Create default role for the resource
	defaultRoleDTO := role.CreateRoleDTO{
		OrgId:    orgID,
		RoleName: "Owner of " + name,
	}
	roleId, err := s.roleRepo.CreateRole(ctx, defaultRoleDTO)
	if err != nil {
		return err
	}
	// add default permission to role
	defaultPermisionDTO := role.ModifyPermissionDTO{
		PermissionName: "Owner of " + name,
		Action:         enum.PermissionActionsWrite,
		ResourceId: resourceID,
	}
	roleId, err = s.roleRepo.AddPermission(ctx, defaultPermisionDTO, roleId)
	if err != nil {
		return err
	}
	// add role to user
	_, err = s.userRepo.AddRole(ctx, userId, roleId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateResource(ctx context.Context, dto CreateResourceDTO, userID string) (string, error) {
	resourceId, err := s.repo.CreateResource(ctx, dto)
	if err != nil {
		fmt.Println("Error creating resource:", err)
		return "", err
	}
	
	// If creating org no need to create rela
	if dto.ParentId == "" {
		err := s.AddDefaultRoleToUser(ctx, userID, resourceId, resourceId, dto.ResourceName)
		if err != nil {
			fmt.Println("Error adding default role to user:", err)
			return "", err
		}
		return resourceId, nil
	}
	pID, err := bson.ObjectIDFromHex(dto.ParentId)
	if err != nil {
		fmt.Println("Invalid parent ID:", err)
		return "", err
	}
	filter := map[string]any{
		"_id": pID,
	}
	parentResource, err := s.repo.GetResourceByFilter(ctx, filter)
	if err != nil {
		fmt.Println("Error getting parent resource:", err)
		return "", err
	}
	// if creating project space use orgId as dto.
	if( parentResource.ResourceType == enum.ResourceTypeOrganization ){
		err = s.AddDefaultRoleToUser(ctx, userID, dto.ParentId, resourceId, dto.ResourceName)
		if err != nil {
			fmt.Println("Error adding default role to user:", err)
			return "", err
		}
	} else {
		// get orgId from project space id (parentID) first 
		orgId, err := s.resourceRelaRepo.GetParentIdByChildId(ctx, dto.ParentId)
		if err != nil {
			fmt.Println("Error getting orgId from project space id:", err)
			return "", err
		}
		err = s.AddDefaultRoleToUser(ctx, userID, orgId, resourceId, dto.ResourceName) 
		if err != nil {
			fmt.Println("Error adding default role to user:", err)
			return "", err
		}
	}
	// create resource relationship
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
