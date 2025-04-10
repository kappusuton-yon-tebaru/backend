package user

import (
	"context"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	repo   *Repository
	logger *logger.Logger
}

func NewService(repo *Repository, logger *logger.Logger) *Service {
	return &Service{
		repo,
		logger,
	}
}

func (s *Service) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) DeleteUserById(ctx context.Context, id string) *werror.WError {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid user id")
	}

	filter := map[string]any{
		"_id": objId,
	}

	count, err := s.repo.DeleteUser(ctx, filter)
	if err != nil {
		return werror.NewFromError(err)
	}

	if count == 0 {
		return werror.New().
			SetCode(http.StatusNotFound).
			SetMessage("User not found")
	}

	return nil
}

func (s *Service) AddRole(ctx context.Context, userId string, roleId string) (string, error) {
	userId, err := s.repo.AddRole(ctx, userId, roleId)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (s *Service) RemoveRole(ctx context.Context, userId string, roleId string) *werror.WError {
	count, err := s.repo.RemoveRole(ctx, userId, roleId)
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
