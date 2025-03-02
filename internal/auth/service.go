package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/user"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type Service struct {
	config      *config.Config
	sessionRepo *Repository
	userRepo    *user.Repository
	logger      *logger.Logger
}

func NewService(config *config.Config, sessionRepo *Repository, repo *user.Repository, logger *logger.Logger) *Service {
	return &Service{
		config,
		sessionRepo,
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

	_, err = s.userRepo.CreateUser(ctx, dto)
	if err != nil && mongo.IsDuplicateKeyError(err) {
		return werror.NewFromError(err).SetMessage("user already exists").SetCode(400)
	} else if err != nil {
		s.logger.Error("error occured while creating user", zap.Error(err))
		return werror.NewFromError(err).SetMessage("error occured while registering")
	}

	return nil
}

func (s *Service) Login(ctx context.Context, dto user.UserCredentialDTO) (string, *werror.WError) {
	user, err := s.userRepo.GetUserByEmail(ctx, dto.Email)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", werror.New().SetMessage("email and password does not match").SetCode(http.StatusUnauthorized)
	} else if err != nil {
		return "", werror.NewFromError(err)
	}

	if !utils.ComparePassword(user.Password, dto.Password) {
		return "", werror.New().SetMessage("email and password does not match").SetCode(http.StatusUnauthorized)
	}

	sessionId, werr := s.CreateSession(ctx, user.Id)
	if werr != nil {
		return "", werr.SetMessage("error occured while logging in")
	}

	token, err := utils.CreateJwtToken(sessionId, s.config.SessionExpiresInSecond, s.config.JwtSecret)
	if err != nil {
		s.logger.Error("error occured while creating jwt token", zap.Error(err))
		return "", werror.NewFromError(err).SetMessage("error occured while logging in")
	}

	return token, nil
}

func (s *Service) CreateSession(ctx context.Context, userId string) (string, *werror.WError) {
	id, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return "", werror.NewFromError(err).SetMessage("invalid user id")
	}

	dto := CreateSessionDTO{
		UserId:    id,
		ExpiresAt: time.Now().Add(time.Duration(s.config.SessionExpiresInSecond) * time.Second),
	}

	sessionId, err := s.sessionRepo.CreateSession(ctx, dto)
	if err != nil {
		s.logger.Error("error occured while creating session", zap.Error(err))
		return "", werror.NewFromError(err)
	}

	return sessionId, nil
}

func (s *Service) GetSession(ctx context.Context, sessionId string) (models.Session, *werror.WError) {
	id, err := bson.ObjectIDFromHex(sessionId)
	if err != nil {
		return models.Session{}, werror.NewFromError(err).SetMessage("invalid session id").SetCode(http.StatusUnauthorized)
	}

	session, err := s.sessionRepo.GetSessionById(ctx, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.Session{}, werror.New().SetMessage("session not found").SetCode(http.StatusNotFound)
	} else if err != nil {
		return models.Session{}, werror.NewFromError(err)
	}

	return session, nil
}
