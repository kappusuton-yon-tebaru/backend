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

func (s *Service) AddUserToUserGroup(ctx context.Context, dto AddUserToUserGroupDTO, id string) (any, error) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}
	addId, err := s.repo.AddUserToUserGroup(ctx, dto, objId)
	if err != nil {
		return "", err
	}

	return addId, nil
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

func (s *Service) DeleteUserFromUserGroupById(ctx context.Context, groupId string, userId string) *werror.WError {
	gId, err := bson.ObjectIDFromHex(groupId)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid user group id")
	}

	uId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid user id")
	}

	filter := map[string]any{
		"group_id": gId,
		"user_id":  uId,
	}

	count, err := s.repo.DeleteUserFromUserGroup(ctx, filter)
	if err != nil {
		return werror.NewFromError(err)
	}

	if count == 0 {
		return werror.New().
			SetCode(http.StatusNotFound).
			SetMessage("User not in group")
	}

	return nil
}
