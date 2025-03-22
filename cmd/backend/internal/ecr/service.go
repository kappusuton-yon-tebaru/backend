package ecr

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
)

type PaginatedECRImages = models.Paginated[ECRImageResponse]

type Service struct {
	repo *ECRRepository
}

func NewService(repo *ECRRepository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) GetECRImages(ctx context.Context, registry models.RegistryProviders, serviceName string, queryParam query.QueryParam) (PaginatedECRImages, error) {
	images, err := s.repo.GetImages(ctx, registry, serviceName, queryParam)
	if err != nil {
		return PaginatedECRImages{}, err
	}

	return images, nil
}
