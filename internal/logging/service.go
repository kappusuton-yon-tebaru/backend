package logging

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) BatchInsertLog(ctx context.Context, dtos []InsertLogDTO) *werror.WError {
	err := s.repo.BatchInsertLog(ctx, dtos)
	if err != nil {
		return werror.NewFromError(err)
	}

	return nil
}
