package role

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

func (s *Service) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	roles, err := s.repo.GetAllRoles(ctx)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *Service) CreateRole(ctx context.Context, dto CreateRoleDTO) (string, error) {
	id, err := s.repo.CreateRole(ctx, dto)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) UpdateRole(ctx context.Context, dto UpdateRoleDTO, id string) (string, *werror.WError) {
	roleId, err := s.repo.UpdateRole(ctx, dto, id)
	if err != nil {
		return "", werror.NewFromError(err).
			SetCode(http.StatusBadRequest)
	}

	return roleId, nil
}

func (s *Service) DeleteRoleById(ctx context.Context, id string) *werror.WError {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid image id")
	}

	filter := map[string]any{
		"_id": objId,
	}

	count, err := s.repo.DeleteRole(ctx, filter)
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

func (s *Service) AddPermission(ctx context.Context, dto ModifyPermissionDTO, id string) (string, *werror.WError) {
	roleId, err := s.repo.AddPermission(ctx, dto, id)
	if err != nil {
		return "", werror.NewFromError(err).
			SetCode(http.StatusBadRequest)
	}

	return roleId, nil
}

func (s *Service) UpdatePermission(ctx context.Context, dto ModifyPermissionDTO, roleId string, permId string) (string, *werror.WError) {
	roleId, err := s.repo.UpdatePermission(ctx, dto, roleId, permId)
	if err != nil {
		return "", werror.NewFromError(err).
			SetCode(http.StatusBadRequest)
	}

	return roleId, nil
}

func (s *Service) DeletePermission(ctx context.Context, roleId string, permId string) *werror.WError {
	count, err := s.repo.DeletePermission(ctx, roleId, permId)
	if err != nil {
		return werror.NewFromError(err).
		SetCode(http.StatusBadRequest)
	}

	if count == 0 {
		return werror.New().
			SetCode(http.StatusNotFound).
			SetMessage("not found")
	}

	return nil
}
