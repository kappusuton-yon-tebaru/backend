package ecr

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
)

type Service struct {
	repo *ECRRepository
}

func NewService(repo *ECRRepository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) GetECRImages(repoURI string, serviceName string, pagination models.Pagination) (models.Paginated[ECRImageResponse], error) {
	images, err := s.repo.GetImages(repoURI, pagination)
	if err != nil {
		return models.Paginated[ECRImageResponse]{}, err
	}

	return images, nil
}
