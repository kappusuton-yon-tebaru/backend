package image

import (
	"context"
	"errors"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
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

func (s *Service) DeleteImage(ctx context.Context, id string) error {
	filter, err := NewFilter().SetId(id).Build()
	if err != nil {
		return err
	}

	count, err := s.repo.DeleteImage(ctx, filter)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("image not found")
	}

	return nil
}
