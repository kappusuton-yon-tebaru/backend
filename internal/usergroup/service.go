package usergroup

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

func (s *Service) GetAllUserGroups(ctx context.Context) ([]models.UserGroup, error) {
	usergroups, err := s.repo.GetAllUserGroups(ctx)
	if err != nil {
		return nil, err
	}

	return usergroups, nil
}

func (s *Service) CreateUserGroup(ctx context.Context, dto CreateUserGroupDTO) (any, error) {
	id, err := s.repo.CreateUserGroup(ctx, dto)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) DeleteUserGroupById(ctx context.Context, id string) *werror.WError {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid user group id")
	}

	filter := map[string]any{
		"_id": objId,
	}

	count, err := s.repo.DeleteUserGroup(ctx, filter)
	if err != nil {
		return werror.NewFromError(err)
	}

	if count == 0 {
		return werror.New().
			SetCode(http.StatusNotFound).
			SetMessage("User group not found")
	}

	return nil
}
