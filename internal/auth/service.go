package auth

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/user"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type Service struct {
	repo   *user.Repository
	logger *logger.Logger
}

func NewService(repo *user.Repository, logger *logger.Logger) *Service {
	return &Service{
		repo,
		logger,
	}
}

func (s *Service) Register(ctx context.Context, dto user.UserCredentialDTO) *werror.WError {
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
