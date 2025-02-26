package ecr

import (
	"strings"

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
	images, err := s.repo.GetImages(repoURI)
	if err != nil {
		return models.Paginated[ECRImageResponse]{}, err
	}

	var response []ECRImageResponse
	for _, img := range images {
		if strings.Contains(img, serviceName) {
			response = append(response, ECRImageResponse{
				ImageTag: img,
			})
		}
	}

	return models.Paginated[ECRImageResponse]{
		Page: pagination.Page,
		Limit: pagination.Limit,
		Total: len(response),
		Data: response,
	} , nil
}
