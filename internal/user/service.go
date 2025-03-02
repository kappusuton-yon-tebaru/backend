package user

import (
	"context"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
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

func (s *Service) Register(ctx context.Context, dto RegisterDTO) *werror.WError {
	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		s.logger.Error("error occured while hashing password", zap.Error(err))
		return werror.NewFromError(err).SetMessage("error occured while registering")
	}

	dto.Password = hashedPassword

	_, err = s.repo.CreateUser(ctx, dto)
	if err != nil && mongo.IsDuplicateKeyError(err) {
		return werror.NewFromError(err).SetMessage("user already exists").SetCode(400)
	} else if err != nil {
		s.logger.Error("error occured while creating user", zap.Error(err))
		return werror.NewFromError(err).SetMessage("error occured while registering")
	}

	return nil
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
