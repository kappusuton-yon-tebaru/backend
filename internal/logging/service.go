package logging

import (
	"context"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.uber.org/zap"
	"k8s.io/kube-openapi/pkg/validation/strfmt/bson"
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

func (s *Service) BatchInsertLog(ctx context.Context, dtos []InsertLogDTO) *werror.WError {
	err := s.repo.BatchInsertLog(ctx, dtos)
	if err != nil {
		return werror.NewFromError(err)
	}

	return nil
}

func (s *Service) GetLog(ctx context.Context, queryParams query.QueryParam) ([]models.Log, *werror.WError) {
	if queryParams.CursorPagination != nil && queryParams.CursorPagination.Cursor != nil {
		_, err := bson.ObjectIDFromHex(*queryParams.CursorPagination.Cursor)
		if err != nil {
			return nil, werror.NewFromError(err).SetMessage("invalid cursor").SetCode(http.StatusBadRequest)
		}
	}

	logs, err := s.repo.GetLog(ctx, queryParams)
	if err != nil {
		s.logger.Error("error occured while getting log", zap.Error(err))
		return nil, werror.NewFromError(err)
	}

	return logs, nil
}
