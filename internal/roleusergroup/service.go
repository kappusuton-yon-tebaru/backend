package roleusergroup

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

func (s *Service) GetAllRoleUserGroups(ctx context.Context) ([]models.RoleUserGroup, error) {
	roleUserGroups, err := s.repo.GetAllRoleUserGroups(ctx)
	if err != nil {
		return nil, err
	}

	return roleUserGroups, nil
}

func (s *Service) CreateRoleUserGroup(ctx context.Context, dto CreateRoleUserGroupDTO) (string, error) {
	id, err := s.repo.CreateRoleUserGroup(ctx, dto)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) DeleteRoleUserGroupById(ctx context.Context, id string) *werror.WError {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid image id")
	}

	filter := map[string]any{
		"_id": objId,
	}

	count, err := s.repo.DeleteRoleUserGroup(ctx, filter)
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