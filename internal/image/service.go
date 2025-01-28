package image

import (
	"context"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
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

func (s *Service) GetAllImages(ctx context.Context) ([]models.Image, error) {
	images, err := s.repo.GetAllImages(ctx)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (s *Service) DeleteImage(ctx context.Context, id string) *werror.WError {
	filter, err := NewFilter().SetId(id).Build()
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid image id")
	}

	count, err := s.repo.DeleteImage(ctx, filter)
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
